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

type JobClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewJobClient(baseURL string) *JobClient {
	return &JobClient{baseURL: baseURL, httpClient: &http.Client{}}
}

func (c *JobClient) CreateJob(ctx context.Context, req ports.CreateJobRequest) (ports.CreateJobResponse, error) {
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

	if auth, ok := ctx.Value("Authorization").(string); ok && auth != "" {
		httpReq.Header.Set("Authorization", auth)
	}

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

func (c *JobClient) GetJobOutcome(ctx context.Context, jobID string) (ports.JobOutcomeResponse, error) {
	url := fmt.Sprintf("%s/jobs/%s/outcome", c.baseURL, jobID)

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return ports.JobOutcomeResponse{}, err
	}

	if auth, ok := ctx.Value("Authorization").(string); ok && auth != "" {
		httpReq.Header.Set("Authorization", auth)
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
