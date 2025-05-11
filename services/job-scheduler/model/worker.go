package model

import (
	"fmt"
	"net/url"
)

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

// -------------------------- Endpoints --------------------------

func PutWorkerStatusEndpoint(base string, id string) string {
	// FIXME: change string to UUID
	return fmt.Sprintf("%s/workers/%s/status", base, id)
}

func GetWorkersEndpoint(base string) string {
	baseUrl := fmt.Sprintf("%s/workers", base)

	params := url.Values{}
	params.Add("status", string(WorkerStatusAvailable))

	fullUrl := baseUrl + "?" + params.Encode()
	return fullUrl
}

// -------------------------- Response & Request --------------------------

type GetWorkersResponse []Worker

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
