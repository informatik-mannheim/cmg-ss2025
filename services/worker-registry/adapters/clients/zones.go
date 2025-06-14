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

	go func() {
		ctx := context.Background()
		sleepStr := os.Getenv("CARBON_INTENSITY_PROVIDER_INTERVAL")
		sleepSec := 60

		if sleepStr != "" {
			if parsed, err := strconv.Atoi(sleepStr); err == nil {
				sleepSec = parsed
			} else {
				fmt.Println("Invalid sleep interval, using default 60s:", err)
			}
		}

		for {
			resp, err := client.GetZones(ctx)
			if err == nil {
				client.zones = resp.Zones
				fmt.Println("Zones successfully loaded:", client.zones)
				return
			}
			fmt.Println("Fetching zones failed...", err)
			time.Sleep(time.Duration(sleepSec) * time.Second)
		}
	}()

	return client
}

func (c *ZoneClient) GetZones(ctx context.Context) (ports.ZoneResponse, error) {
	url := fmt.Sprintf("%s/carbon-intensity/zones", c.baseURL)
	httpReq, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return ports.ZoneResponse{}, fmt.Errorf("failed to create request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ports.ZoneResponse{}, fmt.Errorf("worker-registry error: %s", resp.Status)
	}

	var out ports.ZoneResponse
	err = json.NewDecoder(resp.Body).Decode(&out)
	return out, err
}

func (z *ZoneClient) IsValidZone(code string, ctx context.Context) bool {
	for _, zone := range z.zones {
		if zone.Code == code {
			return true
		}
	}
	return false
}
