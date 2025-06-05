package client_http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/informatik-mannheim/cmg-ss2025/pkg/logging"
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

	logging.From(ctx).Debug("Sending worker registration", "zone", req.Zone, "url", url)

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		logging.From(ctx).Error("Failed to create registration request", "error", err)
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		logging.From(ctx).Error("HTTP request failed during worker registration", "error", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		logging.From(ctx).Warn("Unexpected response during registration", "status", resp.StatusCode, "response", string(body))
		return nil, fmt.Errorf("register worker failed: %s", body)
	}

	var regResp ports.RegisterRespose
	if err := json.NewDecoder(resp.Body).Decode(&regResp); err != nil {
		logging.From(ctx).Error("Failed to decode registration response", "error", err)
		return nil, err
	}

	logging.From(ctx).Debug("Worker registered", "workerID", regResp.ID, "zone", regResp.Zone, "status", regResp.Status)
	return &regResp, nil
}

func (c *RegistryClient) UpdateWorkerStatus(ctx context.Context, req ports.HeartbeatRequest) error {
	url := fmt.Sprintf("%s/workers/%s/status", c.BaseURL, req.WorkerID)

	payload := map[string]string{"workerStatus": req.Status}
	body, err := json.Marshal(payload)
	if err != nil {
		logging.From(ctx).Error("Failed to marshal worker status payload", "workerID", req.WorkerID, "error", err)
		return err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewReader(body))
	if err != nil {
		logging.From(ctx).Error("Failed to create status update request", "workerID", req.WorkerID, "error", err)
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	logging.From(ctx).Debug("Updating worker status", "workerID", req.WorkerID, "status", req.Status)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		logging.From(ctx).Error("HTTP request failed during status update", "workerID", req.WorkerID, "error", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		logging.From(ctx).Warn("Unexpected response during status update", "workerID", req.WorkerID, "status", resp.StatusCode, "response", string(respBody))
		return fmt.Errorf("update worker status failed: %s", respBody)
	}

	logging.From(ctx).Debug("Worker status updated", "workerID", req.WorkerID, "status", req.Status)
	return nil
}
