package ports

import "time"

type JobStatus string

const (
	StatusQueued    JobStatus = "queued"    // default value for new job
	StatusScheduled JobStatus = "scheduled" // is set as soon as a worker has been assigned (by the job-scheduler)
	StatusRunning   JobStatus = "running"   // set by daemon
	StatusCompleted JobStatus = "completed" // set by daemon
	StatusFailed    JobStatus = "failed"    // set by daemon
	StatusCancelled JobStatus = "cancelled" // set by daemon
)

type ContainerImage struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// Json tags helps by (de-)serializing, json:"id" -> "id":"1234", functionality imported by "encoding/json" package

type Job struct {

	// set by job-service, theyre set automatically
	Id        string    `json:"id"`        // generated as UUID
	UserID    string    `json:"userId"`    // get from JWT
	CreatedAt time.Time `json:"createdAt"` // set at creation
	UpdatedAt time.Time `json:"updatedAt"` // set at creation

	// set by consumer-cli, theyre not empty by default
	JobName              string            `json:"jobName"` // set by User
	Image                ContainerImage    `json:"image"`
	AdjustmentParameters map[string]string `json:"parameters"`   // e.g key(-p) : value (8080:8080)
	CreationZone         string            `json:"creationZone"` // origin of the job creation

	// set by job-scheduler
	WorkerID        string `json:"workerId"`        // default value is empty string - saved as UUID
	ComputeZone     string `json:"computeZone"`     // default value is empty string - saved as "zone key", we get from Electricity Maps API, e.g "DE" (germany)
	CarbonIntensity int    `json:"carbonIntensity"` // default value is -1 - CO2eq/kWh which are emitted during job execution
	CarbonSaving    int    `json:"carbonSavings"`   // default value is -1 - consumption savings compared to the actual consumer location

	// set by worker
	Result       string `json:"result"`       // empty string by default - perhaps some containers will provide a result
	ErrorMessage string `json:"errorMessage"` // empty string by default

	// multiple access
	Status JobStatus `json:"status"` // default value is "queued"
}
