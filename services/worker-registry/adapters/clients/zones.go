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
	baseURL    string
	httpClient *http.Client
}

func NewZoneClient(baseURL string) *ZoneClient {
	return &ZoneClient{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}
}

func (c *ZoneClient) GetZones(ctx context.Context) (ports.ZoneResponse, error) {
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
	err = json.NewDecoder(resp.Body).Decode(&out)
	return out, err
}

func (c *ZoneClient) IsValidZone(code string, ctx context.Context) bool {
	zr, err := c.GetZones(ctx)
	if err != nil {
		logging.Warn("Failed to fetch zones for validation: " + err.Error())
		return false
	}
	for _, zone := range zr.Zones {
		if zone.Code == code {
			return true
		}
	}
	return false
}
