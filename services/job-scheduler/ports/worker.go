package ports

import (
	"github.com/google/uuid"
)

type WorkerStatus string

const (
	WorkerStatusAvailable WorkerStatus = "AVAILABLE" // default value for new worker
	WorkerStatusRunning   WorkerStatus = "RUNNING"   // set by Job Scheduler
)

type Worker struct {
	Id     uuid.UUID    `json:"id"`
	Status WorkerStatus `json:"status"`
	Zone   string       `json:"zone"`
}

type GetWorkersResponse []Worker

// This struct is used for the patch-request to the worker service
type UpdateWorkerPayload struct {
	WorkerStatus WorkerStatus `json:"status"` // default (and probably only) value is "running"
}

// This struct is returned by the worker service as response to the put-request
type UpdateWorkerResponse GetWorkersResponse

type UpdateWorker struct {
	ID uuid.UUID `json:"id"`
	// No status on this struct, because there is only 1 possible option,
	// so the function will set it itself.
}

type WorkerAdapter interface {
	GetWorkers() (GetWorkersResponse, error)
	AssignWorker(update UpdateWorker) error
}
