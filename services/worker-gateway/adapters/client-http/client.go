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

type HTTPClient struct {
	JobServiceURL      string
	RegistryServiceURL string
	httpClient         *http.Client
}

func NewHTTPClient(jobURL, registryURL string) *HTTPClient {
	return &HTTPClient{
		JobServiceURL:      jobURL,
		RegistryServiceURL: registryURL,
		httpClient:         &http.Client{},
	}
}

func (c *HTTPClient) addHeaders(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
}

func (c *HTTPClient) RegisterWorker(ctx context.Context, req ports.RegisterRequest) error {
	url := fmt.Sprintf("%s/workers?zone=%s", c.RegistryServiceURL, url.QueryEscape(req.Location))

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return err
	}
	c.addHeaders(httpReq)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("register worker failed: %s", body)
	}

	return nil
}

func (c *HTTPClient) UpdateWorkerStatus(ctx context.Context, req ports.HeartbeatRequest) error {
	url := fmt.Sprintf("%s/workers/%s/status", c.RegistryServiceURL, req.WorkerID)

	payload := map[string]string{"workerStatus": req.Status}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	c.addHeaders(httpReq)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("update worker status failed: %s", respBody)
	}

	return nil
}

func (c *HTTPClient) UpdateJob(ctx context.Context, req ports.ResultRequest) error {
	url := fmt.Sprintf("%s/jobs/%s/update-workerdaemon", c.JobServiceURL, req.JobID)

	payload := map[string]string{
		"status":       req.Status,
		"result":       req.Result,
		"errorMessage": req.ErrorMessage,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPatch, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	c.addHeaders(httpReq)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("update job failed: %s", respBody)
	}

	return nil
}

func (c *HTTPClient) FetchScheduledJobs(ctx context.Context) ([]ports.Job, error) {
	url := fmt.Sprintf("%s/jobs?status=scheduled", c.JobServiceURL)

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	c.addHeaders(httpReq)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return []ports.Job{}, nil
	}

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("fetch scheduled jobs failed: %s", respBody)
	}

	var jobs []ports.Job
	if err := json.NewDecoder(resp.Body).Decode(&jobs); err != nil {
		return nil, err
	}

	return jobs, nil
}
