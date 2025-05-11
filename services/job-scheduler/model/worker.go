package model

import "fmt"

type WorkerStatus string

const (
	WorkerStatusAvailable WorkerStatus = "AVAILABLE" // default value for new worker
	WorkerStatusRunning   WorkerStatus = "RUNNING"   // set by Job Scheduler
)

type Worker struct {
	Id     string       `json:"id"` // FIXME: actually UUID
	Status WorkerStatus `json:"status"`
	Zone   string       `json:"zone"`
}

func PutWorkerStatusEndpoint(id string) string {
	// FIXME: Add base
	// FIXME: change string to UUID
	return fmt.Sprintf("TODO:ADDBASE/workers/%s/status", id)
}

// This struct is used for the patch-request to the worker service
type UpdateWorkerPayload struct {
	WorkerStatus WorkerStatus `json:"workerStatus"` // default (and probably only) value is "running"
}

// This struct is returned by the worker service as response to the put-request
type UpdateWorkerResponse []struct {
	ID     string       `json:"id"` // generated as UUID
	Status WorkerStatus `json:"status"`
	Zone   string       `json:"zone"`
}
