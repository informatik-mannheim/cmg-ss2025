package ports

import "github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/model"

// FIXME: Change id-string to id-uuid, currently not possible due to lack of
// uuid package approval...

type UpdateJob struct {
	ID              string  `json:"id"`       // actually uuid
	WorkerID        string  `json:"workerId"` // actually uuid
	ComputeZone     string  `json:"computeZone"`
	CarbonIntensity float64 `json:"carbonIntensity"`
	CarbonSavings   float64 `json:"carbonSavings"`
	// No status on this struct, because there is only 1 possible option,
	// so the function will set it itself.
}

type UpdateWorker struct {
	ID string `json:"id"` // actually uuid
	// No status on this struct, because there is only 1 possible option,
	// so the function will set it itself.
}

type Notifier interface {
	// -- Getters --
	GetJobs() ([]model.Job, error)
	GetWorkers() ([]model.Worker, error)
	GetCarbonIntensities(zones []string) ([]model.CarbonIntensityResponse, error)

	// -- Setters --
	AssignJob(update UpdateJob) error
	AssignWorker(update UpdateWorker) error
}
