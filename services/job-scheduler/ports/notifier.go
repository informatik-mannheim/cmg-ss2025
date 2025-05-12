package ports

import "github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/model"

type UpdateJob struct {
	ID              string  `json:"id"`              // FIXME: actually UUID
	WorkerID        string  `json:"workerId"`        // FIXME: actually UUID
	ComputeZone     string  `json:"computeZone"`     //
	CarbonIntensity float64 `json:"carbonIntensity"` //
	CarbonSavings   float64 `json:"carbonSavings"`   //
	// No status on this struct, because there is only 1 possible option,
	// so the function will set it itself.
}

type UpdateWorker struct {
	ID string `json:"id"` // FIXME: actually UUID
	// No status on this struct, because there is only 1 possible option,
	// so the function will set it itself.
}

type Notifier interface {
	// -- Getters --
	GetJobs() (model.GetJobsResponse, error)
	GetWorkers() (model.GetWorkersResponse, error)
	GetCarbonIntensities(zones []string) (model.CarbonIntensityResponse, error)

	// -- Setters --
	AssignJob(update UpdateJob) error
	AssignWorker(update UpdateWorker) error
}
