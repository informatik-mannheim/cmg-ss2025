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
	Name    string `json:"name" db:"image_name"`
	Version string `json:"version" db:"image_version"`
}

// Json tags helps by (de-)serializing, json:"id" -> "id":"1234", functionality imported by "encoding/json" package

type Job struct {

	// set by job-service, theyre set automatically
	Id        string    `json:"id" db:"id"`                // generated as UUID
	UserID    string    `json:"userId" db:"user_id"`       // get from JWT
	CreatedAt time.Time `json:"createdAt" db:"created_at"` // set at creation
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"` // set at creation

	// set by consumer-cli, theyre not empty by default
	JobName              string            `json:"jobName" db:"job_name"` // set by User
	Image                ContainerImage    `json:"image" db:"-"`
	AdjustmentParameters map[string]string `json:"parameters" db:"adjustment_parameters"` // e.g key(-p) : value (8080:8080)
	CreationZone         string            `json:"creationZone" db:"creation_zone"`       // origin of the job creation

	// set by job-scheduler
	WorkerID        string `json:"workerId" db:"worker_id"`               // default value is empty string - saved as UUID
	ComputeZone     string `json:"computeZone" db:"compute_zone"`         // default value is empty string - saved as "zone key", we get from Electricity Maps API, e.g "DE" (germany)
	CarbonIntensity int    `json:"carbonIntensity" db:"carbon_intensity"` // default value is -1 - CO2eq/kWh which are emitted during job execution
	CarbonSaving    int    `json:"carbonSavings" db:"carbon_savings"`     // default value is -1 - consumption savings compared to the actual consumer location

	// set by worker
	Result       string `json:"result" db:"result"`              // empty string by default - perhaps some containers will provide a result
	ErrorMessage string `json:"errorMessage" db:"error_message"` // empty string by default

	// multiple access
	Status JobStatus `json:"status" db:"job_status"` // default value is "queued"
}
