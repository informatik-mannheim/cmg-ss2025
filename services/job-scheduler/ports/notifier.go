package ports

import (
	CarbonIntensityProvider "github.com/informatik-mannheim/cmg-ss2025/services/carbon-intensity-provider/model"
	Job "github.com/informatik-mannheim/cmg-ss2025/services/job"
	WorkerRegistry "github.com/informatik-mannheim/cmg-ss2025/services/worker-registry/ports"
)

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
	GetJobs() ([]Job.Job, error)
	GetWorkers() ([]WorkerRegistry.Worker, error)
	GetCarbonIntensities(zones []string) ([]CarbonIntensityProvider.CarbonIntensityData, error)

	// -- Setters --
	AssignJob(update UpdateJob) error
	AssignWorker(update UpdateWorker) error
}
