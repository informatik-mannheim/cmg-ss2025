package client_http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/informatik-mannheim/cmg-ss2025/services/worker-gateway/ports"
)

type JobClient struct {
	BaseURL    string
	httpClient *http.Client
	//testJobs   []ports.Job
}

func NewJobClient(baseURL string) *JobClient {
	return &JobClient{BaseURL: baseURL, httpClient: &http.Client{}}
	/*testJobs: []ports.Job{
		{
			ID:       "job123",
			WorkerID: "worker123",
			Status:   "scheduled",
			Result:   "",
			ErrorMsg: "",
		},
		{
			ID:       "job456",
			WorkerID: "worker123",
			Status:   "scheduled",
			Result:   "",
			ErrorMsg: "",
		},
		{
			ID:       "job789",
			WorkerID: "worker123",
			Status:   "scheduled",
			Result:   "",
			ErrorMsg: "",
		},
	},
	*/

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
		return err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPatch, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return err
		//filtered := make([]ports.Job, 0, len(c.testJobs))
		//for _, job := range c.testJobs {
		//	if job.ID != req.JobID {
		//		filtered = append(filtered, job)
		//	}
		//}
		//c.testJobs = filtered
		//return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("update job failed: %s", respBody)
	}

	return nil
}

func (c *JobClient) FetchScheduledJobs(ctx context.Context) ([]ports.Job, error) {
	url := fmt.Sprintf("%s/jobs?status=scheduled", c.BaseURL)

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
		//var testJobs1 = c.testJobs
		//return testJobs1, nil
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
