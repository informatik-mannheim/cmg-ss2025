package client_http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/informatik-mannheim/cmg-ss2025/pkg/logging"
	"github.com/informatik-mannheim/cmg-ss2025/services/worker-gateway/ports"
)

type JobClient struct {
	BaseURL    string
	httpClient *http.Client
}

func NewJobClient(baseURL string) *JobClient {
	return &JobClient{BaseURL: baseURL, httpClient: &http.Client{}}
}

func (c *JobClient) UpdateJob(ctx context.Context, req ports.ResultRequest) error {
	url := fmt.Sprintf("%s/jobs/%s/update-workerdaemon", c.BaseURL, req.JobID)

	payload := map[string]string{
		"status":       req.Status,
		"result":       req.Result,
		"errorMessage": req.ErrorMessage,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		logging.From(ctx).Error("Failed to marshal job update payload", "jobID", req.JobID, "error", err)
		return err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPatch, url, bytes.NewReader(body))
	if err != nil {
		logging.From(ctx).Error("Failed to create job update request", "jobID", req.JobID, "error", err)
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	logging.From(ctx).Debug("Sending job update", "jobID", req.JobID, "status", req.Status)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		logging.From(ctx).Error("HTTP request failed during job update", "jobID", req.JobID, "error", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		respBody, _ := io.ReadAll(resp.Body)
		logging.From(ctx).Warn("Unexpected response during job update", "jobID", req.JobID, "status", resp.StatusCode, "response", string(respBody))
		return fmt.Errorf("update job failed: %s", respBody)
	}

	logging.From(ctx).Debug("Job updated successfully", "jobID", req.JobID, "status", req.Status)
	return nil
}

func (c *JobClient) FetchScheduledJobs(ctx context.Context) ([]ports.Job, error) {
	url := fmt.Sprintf("%s/jobs?status=scheduled", c.BaseURL)

	logging.From(ctx).Debug("Fetching scheduled jobs", "url", url)

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		logging.From(ctx).Error("Failed to create request for fetching jobs", "error", err)
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		logging.From(ctx).Error("HTTP request failed during job fetch", "error", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		logging.From(ctx).Debug("No scheduled jobs available")
		return []ports.Job{}, nil
	}

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		logging.From(ctx).Warn("Unexpected response when fetching jobs", "status", resp.StatusCode, "response", string(respBody))
		return nil, fmt.Errorf("fetch scheduled jobs failed: %s", respBody)
	}

	var jobs []ports.Job
	if err := json.NewDecoder(resp.Body).Decode(&jobs); err != nil {
		logging.From(ctx).Error("Failed to decode scheduled jobs response", "error", err)
		return nil, err
	}

	logging.From(ctx).Debug("Scheduled jobs fetched", "count", len(jobs))
	return jobs, nil
}
