package job

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/uuid"
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/ports"
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/utils"
)

func PatchJobStatusEndpoint(base string, id uuid.UUID) string {
	return fmt.Sprintf("%s/jobs/%s/update-scheduler", base, id)
}

func GetJobsEndpoint(base string) string {
	baseUrl := fmt.Sprintf("%s/jobs", base)

	status := []string{string(ports.JobStatusScheduled), string(ports.JobStatusQueued)}

	params := url.Values{}
	params.Add("status", strings.Join(status, ","))

	fullUrl := baseUrl + "?" + params.Encode()
	return fullUrl
}

type JobAdapter struct {
	baseUrl string
	client  http.Client
}

var _ ports.JobAdapter = (*JobAdapter)(nil)

func NewJobAdapter(client http.Client, baseUrl string) *JobAdapter {
	return &JobAdapter{
		baseUrl: baseUrl,
		client:  client,
	}
}

func (adapter *JobAdapter) GetJobs() (ports.GetJobsResponse, error) {
	// For now its kept simple and return an error as soon as it gets one, changes in Phase 3
	endpoint := GetJobsEndpoint(adapter.baseUrl)

	// StatusCode is not relevant yet
	data, _, err := utils.GetRequest[ports.GetJobsResponse](&adapter.client, endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get jobs: %w", err)
	}

	return data, nil
}

func (adapter *JobAdapter) AssignJob(update ports.UpdateJob) error {
	// For now its kept simple and return an error as soon as it gets one, changes in Phase 3
	endpoint := PatchJobStatusEndpoint(adapter.baseUrl, update.ID)

	payload := ports.UpdateJobPayload{
		WorkerID:        update.WorkerID,
		ComputeZone:     update.ComputeZone,
		CarbonIntensity: int(update.CarbonIntensity),
		CarbonSaving:    int(update.CarbonSavings),
		Status:          ports.JobStatusScheduled, // Hardcoded because nothing else is possible
	}

	// StatusCode is not relevant yet
	_, _, err := utils.PatchRequest[ports.UpdateJobPayload, ports.UpdateJobResponse](&adapter.client, endpoint, payload)
	if err != nil {
		return fmt.Errorf("failed to assign job: %w", err)
	}

	return nil
}
