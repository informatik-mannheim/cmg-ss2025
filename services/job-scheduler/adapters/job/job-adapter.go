package job

import (
	"fmt"

	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/model"
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/ports"
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/utils"
)

type JobAdapter struct {
	environments model.Environments
}

var _ ports.JobAdapter = (*JobAdapter)(nil)

func NewJobAdapter(environments model.Environments) *JobAdapter {
	return &JobAdapter{
		environments: environments,
	}
}

func (adapter *JobAdapter) GetJobs() (model.GetJobsResponse, error) {
	// For now its kept simple and return an error as soon as it gets one, changes in Phase 3
	endpoint := model.GetJobsEndpoint(adapter.environments.JobServiceUrl)

	data, err := utils.GetRequest[model.GetJobsResponse](endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get jobs: %w", err)
	}

	return data, nil
}

func (adapter *JobAdapter) AssignJob(update ports.UpdateJob) error {
	// For now its kept simple and return an error as soon as it gets one, changes in Phase 3
	endpoint := model.PatchJobStatusEndpoint(adapter.environments.JobServiceUrl, update.ID)

	payload := model.UpdateJobPayload{
		WorkerID:        update.WorkerID,
		ComputeZone:     update.ComputeZone,
		CarbonIntensity: int(update.CarbonIntensity),
		CarbonSaving:    int(update.CarbonSavings),
		Status:          model.JobStatusScheduled, // Hardcoded because nothing else is possible
	}

	_, err := utils.PatchRequest[model.UpdateJobPayload, model.UpdateJobResponse](endpoint, payload)
	if err != nil {
		return fmt.Errorf("failed to assign job: %w", err)
	}

	return nil
}
