package electricitymaps

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/informatik-mannheim/cmg-ss2025/services/carbon-intensity-provider/ports"
)

// Zone holds metadata for an electricity zone
type Zone struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// Fetcher fetches carbon intensity data for each zone using a dedicated token.
type Fetcher struct {
	TokenByZone map[string]string
	Notifier    ports.Notifier
}

// NewFromEnv creates a new fetcher with tokens loaded from environment variables.
func NewFromEnv(notifier ports.Notifier) *Fetcher {
	tokens := map[string]string{}
	for _, env := range os.Environ() {
		if len(env) > 6 && env[:6] == "TOKEN_" {
			parts := strings.SplitN(env, "=", 2)
			if len(parts) == 2 {
				zone := parts[0][6:]                      // strip TOKEN_
				zone = strings.ReplaceAll(zone, "_", "-") // TOKEN_US_NY → US-NY
				tokens[zone] = parts[1]
			}
		}
	}
	return &Fetcher{
		TokenByZone: tokens,
		Notifier:    notifier,
	}
}

// Fetch retrieves the current carbon intensity data for the given zone.
func (f *Fetcher) Fetch(zone string, ctx context.Context) (ports.CarbonIntensityData, error) {
	token, ok := f.TokenByZone[zone]
	if !ok || token == "" {
		f.Notifier.Event("No token configured for zone: " + zone)
		return ports.CarbonIntensityData{}, fmt.Errorf("no token configured for zone %s", zone)
	}

	url := fmt.Sprintf("https://api.electricitymap.org/v3/carbon-intensity/latest?zone=%s", zone)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		f.Notifier.Event("Failed to create request for zone: " + zone)
		return ports.CarbonIntensityData{}, err
	}
	req.Header.Set("auth-token", token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		f.Notifier.Event("HTTP request failed for zone: " + zone + " — " + err.Error())
		return ports.CarbonIntensityData{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		f.Notifier.Event(fmt.Sprintf("Non-200 response for zone %s: %d", zone, res.StatusCode))
		return ports.CarbonIntensityData{}, fmt.Errorf("API returned status: %d", res.StatusCode)
	}

	var parsed struct {
		CarbonIntensity float64 `json:"carbonIntensity"`
	}

	if err := json.NewDecoder(res.Body).Decode(&parsed); err != nil {
		f.Notifier.Event("Failed to decode response for zone: " + zone)
		return ports.CarbonIntensityData{}, err
	}

	f.Notifier.Event(fmt.Sprintf("Fetched %s with intensity %.2f", zone, parsed.CarbonIntensity))
	return ports.CarbonIntensityData{
		Zone:            zone,
		CarbonIntensity: parsed.CarbonIntensity,
	}, nil
}

// AllElectricityMapZones fetches all available zones from the Electricity Maps API
func (f *Fetcher) AllElectricityMapZones(ctx context.Context) ([]Zone, error) {
	url := "https://api.electricitymap.org/v3/zones"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		f.Notifier.Event("Failed to build zone list request")
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		f.Notifier.Event("Failed to call zone list endpoint")
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		f.Notifier.Event(fmt.Sprintf("Zone list failed: %d", res.StatusCode))
		return nil, fmt.Errorf("failed to fetch zones: %d", res.StatusCode)
	}

	var raw map[string]struct {
		CountryName string `json:"countryName"`
		ZoneName    string `json:"zoneName"`
	}

	if err := json.NewDecoder(res.Body).Decode(&raw); err != nil {
		f.Notifier.Event("Failed to decode zone list")
		return nil, err
	}

	zones := make([]Zone, 0, len(raw))
	for code, meta := range raw {
		name := meta.ZoneName
		if name == "" {
			name = meta.CountryName
		}
		zones = append(zones, Zone{
			Code: code,
			Name: name,
		})
	}
	f.Notifier.Event(fmt.Sprintf("Fetched %d detailed zones (unauthenticated)", len(zones)))
	return zones, nil
}

// GetConfiguredZones returns only zones for which we have tokens configured.
func (f *Fetcher) GetConfiguredZones(ctx context.Context) ([]string, error) {
	configured := make([]string, 0, len(f.TokenByZone))
	for zone := range f.TokenByZone {
		configured = append(configured, zone)
	}
	f.Notifier.Event(fmt.Sprintf("Discovered %d configured zones with tokens", len(configured)))
	return configured, nil
}
