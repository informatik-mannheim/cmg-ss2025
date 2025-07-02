package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/informatik-mannheim/cmg-ss2025/services/worker-registry/ports"
)

type ZoneClient struct {
	baseURL    string
	httpClient *http.Client
	zones      []ports.Zone
}

var _ ports.ZoneClient = (*ZoneClient)(nil)

func NewZoneClient(baseURL string) *ZoneClient {
	client := &ZoneClient{
		baseURL:    baseURL,
		httpClient: &http.Client{},
		zones:      []ports.Zone{},
	}

	go client.refreshZonesLoop()
	return client
}

func (c *ZoneClient) GetZones(ctx context.Context) (ports.ZoneResponse, error) {
	if len(c.zones) == 0 {
		zr, err := c.doGetZones(ctx)
		if err != nil {
			return ports.ZoneResponse{}, err
		}
		c.zones = zr.Zones
		return zr, nil
	}
	return ports.ZoneResponse{Zones: c.zones}, nil
}

func (c *ZoneClient) doGetZones(ctx context.Context) (ports.ZoneResponse, error) {
	url := fmt.Sprintf("%s/carbon-intensity/zones", c.baseURL)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	req.Header.Set("Content-Type", "application/json")

	if auth, ok := ctx.Value("authHeader").(string); ok && auth != "" {
		req.Header.Set("Authorization", auth)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return ports.ZoneResponse{}, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ports.ZoneResponse{}, fmt.Errorf("zone-service error: %s", resp.Status)
	}

	var out ports.ZoneResponse
	err = json.NewDecoder(resp.Body).Decode(&out)
	return out, err
}

func (c *ZoneClient) refreshZonesLoop() {
	sleepSec := 60
	if s := os.Getenv("CARBON_INTENSITY_PROVIDER_INTERVAL"); s != "" {
		if i, err := strconv.Atoi(s); err == nil {
			sleepSec = i
		}
	}
	for {
		if zr, err := c.doGetZones(context.Background()); err == nil {
			c.zones = zr.Zones
			return
		}
		time.Sleep(time.Duration(sleepSec) * time.Second)
	}
}

func (c *ZoneClient) IsValidZone(code string, ctx context.Context) bool {
	for _, zone := range c.zones {
		if zone.Code == code {
			return true
		}
	}
	return false
}
