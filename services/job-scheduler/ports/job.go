package ports

import (
	"github.com/google/uuid"
)

type JobStatus string

const (
	JobStatusQueued    JobStatus = "queued"    // default value for new job
	JobStatusScheduled JobStatus = "scheduled" // is set as soon as a worker has been assigned (by the job-scheduler)
	JobStatusRunning   JobStatus = "running"   // set by daemon
	JobStatusCompleted JobStatus = "completed" // set by daemon
	JobStatusFailed    JobStatus = "failed"    // set by daemon
	JobStatusCancelled JobStatus = "cancelled" // set by daemon
)

type Job struct {

	// set by job-service, theyre set automatically
	ID uuid.UUID `json:"id"` // generated as UUID

	// set by consumer-cli, theyre not empty by default
	CreationZone string `json:"creationZone"` // origin of the job creation

	// set by job-scheduler
	WorkerID        string `json:"workerId"`        // default value is empty string - saved as UUID
	ComputeZone     string `json:"computeZone"`     // default value is empty string - saved as "zone key", we get from Electricity Maps API, e.g "DE" (germany)
	CarbonIntensity int    `json:"carbonIntensity"` // default value is -1 - CO2eq/kWh which are emitted during job execution
	CarbonSaving    int    `json:"carbonSavings"`   // default value is -1 - consumption savings compared to the actual consumer location

	// multiple access
	Status JobStatus `json:"status"` // default value is "queued"
}

// This struct is used for the get-request to the job service
type GetJobsResponse []Job

// This struct is used for the get-request to the job service
type GetJobsError struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// This struct is used for the patch-request to the job service
type UpdateJobPayload struct {
	WorkerID        uuid.UUID `json:"workerId"`        //
	ComputeZone     string    `json:"computeZone"`     //
	CarbonIntensity int       `json:"carbonIntensity"` //
	CarbonSaving    int       `json:"carbonSavings"`   //
	Status          JobStatus `json:"status"`          // default (and probably only) value is "scheduled"
}

// This struct is returned by the job service as response to the patch-request
type UpdateJobResponse struct {
	JobID           uuid.UUID `json:"jobId"`           //
	ComputeZone     string    `json:"computeZone"`     //
	CarbonIntensity int       `json:"carbonIntensity"` //
	CarbonSaving    int       `json:"carbonSavings"`   //
	Status          JobStatus `json:"status"`          // default (and probably only) value is "scheduled"
}

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
	GetJobs() (GetJobsResponse, error)
	AssignJob(update UpdateJob) error
}
