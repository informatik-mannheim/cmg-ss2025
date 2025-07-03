package client_http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/informatik-mannheim/cmg-ss2025/services/consumer-gateway/ports"
)

type ZoneClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewZoneClient(baseURL string) *ZoneClient {
	return &ZoneClient{baseURL: baseURL, httpClient: &http.Client{}}
}

func (c *ZoneClient) GetZone(ctx context.Context, req ports.ZoneRequest) (ports.ZoneResponse, error) {
	url := fmt.Sprintf("%s/zones", c.baseURL)
	body, _ := json.Marshal(req)
	httpReq, _ := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return ports.ZoneResponse{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return ports.ZoneResponse{}, fmt.Errorf("job-service error: %s", resp.Status)
	}

	PingJobScheduler()

	var out ports.ZoneResponse
	err = json.NewDecoder(resp.Body).Decode(&out)
	return out, err
}
