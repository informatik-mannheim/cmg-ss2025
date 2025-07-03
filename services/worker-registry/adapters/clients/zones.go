package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/informatik-mannheim/cmg-ss2025/pkg/logging"
	"github.com/informatik-mannheim/cmg-ss2025/services/worker-registry/ports"
)

type ZoneClient struct {
	baseURL       string
	httpClient    *http.Client
	zones         []ports.Zone
	lastFetched   time.Time
	cacheDuration time.Duration
}

func NewZoneClient(baseURL string, cacheDuration time.Duration) *ZoneClient {
	c := &ZoneClient{
		baseURL:       baseURL,
		httpClient:    &http.Client{Timeout: 5 * time.Second},
		cacheDuration: cacheDuration,
	}

	go func() {
		if _, err := c.GetZones(context.Background()); err != nil {
			logging.Warn("Initial zones fetch failed: " + err.Error())
		}

		ticker := time.NewTicker(c.cacheDuration)
		defer ticker.Stop()

		for range ticker.C {
			if _, err := c.GetZones(context.Background()); err != nil {
				logging.Warn("Periodic zones fetch failed: " + err.Error())
			}
		}
	}()

	return c
}

func (c *ZoneClient) GetZones(ctx context.Context) (ports.ZoneResponse, error) {
	if time.Since(c.lastFetched) > c.cacheDuration {
		zr, err := c.doGetZones(ctx)
		if err != nil {
			logging.Warn("Zones could not be fetched due to error" + err.Error())
			return ports.ZoneResponse{}, err
		}
		c.zones = zr.Zones
		c.lastFetched = time.Now()
	}
	return ports.ZoneResponse{Zones: c.zones}, nil
}

func (c *ZoneClient) doGetZones(ctx context.Context) (ports.ZoneResponse, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/carbon-intensity/zones", nil)
	req.Header.Set("Content-Type", "application/json")
	if auth, ok := ctx.Value("authHeader").(string); ok && auth != "" {
		req.Header.Set("Authorization", auth)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return ports.ZoneResponse{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return ports.ZoneResponse{}, fmt.Errorf("zone-service error: %s", resp.Status)
	}
	var out ports.ZoneResponse
	return out, json.NewDecoder(resp.Body).Decode(&out)
}

func (c *ZoneClient) IsValidZone(code string, ctx context.Context) bool {
	for _, zone := range c.zones {
		if zone.Code == code {
			return true
		}
	}
	return false
}
