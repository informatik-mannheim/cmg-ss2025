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

// Default API URLs (overridable in tests)
var (
	FetchURL        = "https://api.electricitymap.org/v3/carbon-intensity/latest?zone=%s"
	ZoneMetadataURL = "https://api.electricitymap.org/v3/zones"
)

// Fetcher fetches carbon intensity data using tokens per zone
type Fetcher struct {
	TokenByZone map[string]string
	Client      *http.Client
}

// NewFromEnv creates a Fetcher with tokens from environment
func NewFromEnv(notifier ports.Notifier) *Fetcher {
	tokens := map[string]string{}
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, "TOKEN_") {
			parts := strings.SplitN(env, "=", 2)
			if len(parts) == 2 {
				zone := strings.ReplaceAll(parts[0][6:], "_", "-")
				tokens[zone] = parts[1]
			}
		}
	}
	return &Fetcher{
		TokenByZone: tokens,
		Client:      http.DefaultClient,
	}
}

// NewWithClient creates a Fetcher with a custom HTTP client (for testing)
func NewWithClient(notifier ports.Notifier, client *http.Client) *Fetcher {
	return &Fetcher{
		TokenByZone: map[string]string{},
		Client:      client,
	}
}

// Fetch gets the carbon intensity for a specific zone
func (f *Fetcher) Fetch(zone string, ctx context.Context) (ports.CarbonIntensityData, error) {
	token, ok := f.TokenByZone[zone]
	if !ok || token == "" {
		return ports.CarbonIntensityData{}, fmt.Errorf("no token configured for zone %s", zone)
	}

	url := fmt.Sprintf(FetchURL, zone)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return ports.CarbonIntensityData{}, err
	}
	req.Header.Set("auth-token", token)

	res, err := f.Client.Do(req)
	if err != nil {
		return ports.CarbonIntensityData{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return ports.CarbonIntensityData{}, fmt.Errorf("API returned status: %d", res.StatusCode)
	}

	var parsed struct {
		CarbonIntensity float64 `json:"carbonIntensity"`
	}
	if err := json.NewDecoder(res.Body).Decode(&parsed); err != nil {
		return ports.CarbonIntensityData{}, err
	}

	return ports.CarbonIntensityData{
		Zone:            zone,
		CarbonIntensity: parsed.CarbonIntensity,
	}, nil
}

// AllElectricityMapZones returns all zones (unauthenticated)
func (f *Fetcher) AllElectricityMapZones(ctx context.Context) ([]Zone, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ZoneMetadataURL, nil)
	if err != nil {
		return nil, err
	}

	res, err := f.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch zones: %d", res.StatusCode)
	}

	var raw map[string]struct {
		CountryName string `json:"countryName"`
		ZoneName    string `json:"zoneName"`
	}
	if err := json.NewDecoder(res.Body).Decode(&raw); err != nil {
		return nil, err
	}

	zones := make([]Zone, 0, len(raw))
	for code, meta := range raw {
		name := meta.ZoneName
		if name == "" {
			name = meta.CountryName
		}
		zones = append(zones, Zone{Code: code, Name: name})
	}
	return zones, nil
}

// GetConfiguredZones returns zones with a configured token
func (f *Fetcher) GetConfiguredZones(ctx context.Context) ([]string, error) {
	zones := make([]string, 0, len(f.TokenByZone))
	for zone := range f.TokenByZone {
		zones = append(zones, zone)
	}
	return zones, nil
}
