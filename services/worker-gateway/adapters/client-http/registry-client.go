package client_http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/informatik-mannheim/cmg-ss2025/services/worker-gateway/ports"
)

type RegistryClient struct {
	BaseURL    string
	httpClient *http.Client
}

func NewRegistryClient(baseURL string) *RegistryClient {
	return &RegistryClient{BaseURL: baseURL, httpClient: &http.Client{}}
}

func (c *RegistryClient) RegisterWorker(ctx context.Context, req ports.RegisterRequest) (*ports.RegisterRespose, error) {
	url := fmt.Sprintf("%s/workers?zone=%s", c.BaseURL, url.QueryEscape(req.Zone))

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
		//return &ports.RegisterRespose{ID: "worker123", Status: "AVAILABLE", Zone: "DE"}, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("register worker failed: %s", body)
	}

	var regResp ports.RegisterRespose
	if err := json.NewDecoder(resp.Body).Decode(&regResp); err != nil {
		return nil, err
	}

	return &regResp, nil
}

func (c *RegistryClient) UpdateWorkerStatus(ctx context.Context, req ports.HeartbeatRequest) error {
	url := fmt.Sprintf("%s/workers/%s/status", c.BaseURL, req.WorkerID)

	payload := map[string]string{"workerStatus": req.Status}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return err
		//return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("update worker status failed: %s", respBody)
	}

	return nil
}
