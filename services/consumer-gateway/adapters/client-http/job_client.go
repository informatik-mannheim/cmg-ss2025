package client_http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/informatik-mannheim/cmg-ss2025/services/consumer-gateway/ports"
	"io"
	"net/http"
)

type Client struct {
	baseURL string
}

func New(baseURL string) *Client {
	return &Client{baseURL: baseURL}
}

func (c *Client) CreateJob(ctx context.Context, req ports.CreateJobRequest) (ports.CreateJobResponse, error) {
	url := fmt.Sprintf("%s/jobs", c.baseURL)

	body, err := json.Marshal(req)
	if err != nil {
		return ports.CreateJobResponse{}, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return ports.CreateJobResponse{}, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return ports.CreateJobResponse{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return ports.CreateJobResponse{}, fmt.Errorf("job-service error: %s", resp.Status)
	}

	var out ports.CreateJobResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return ports.CreateJobResponse{}, err
	}

	return out, nil
}

func (c *Client) GetJobOutcome(ctx context.Context, jobID string) (ports.JobOutcomeResponse, error) {
	url := fmt.Sprintf("%s/jobs/%s/outcome", c.baseURL, jobID)

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return ports.JobOutcomeResponse{}, err
	}

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return ports.JobOutcomeResponse{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return ports.JobOutcomeResponse{}, fmt.Errorf("job-service error: %s", resp.Status)
	}

	var out ports.JobOutcomeResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return ports.JobOutcomeResponse{}, err
	}

	return out, nil
}
