package ports

import (
	"github.com/google/uuid"
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/model"
)

type UpdateJob struct {
	ID              uuid.UUID `json:"id"`
	WorkerID        uuid.UUID `json:"workerId"`
	ComputeZone     string    `json:"computeZone"`
	CarbonIntensity float64   `json:"carbonIntensity"`
	CarbonSavings   float64   `json:"carbonSavings"`
	// No status on this struct, because there is only 1 possible option,
	// so the function will set it itself.
}

type JobAdapter interface {
	GetJobs() (model.GetJobsResponse, error)
	AssignJob(update UpdateJob) error
}
