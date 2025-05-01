package models

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
	ID                   string            `json:"id"`      // generated as UUID
	JobName              string            `json:"jobName"` //set by User
	UserID               string            `json:"userId"`  // get from JWT
	Image                ContainerImage    `json:"image"`
	AdjustmentParameters map[string]string `json:"parameters"` // e.g key(-p) : value (8080:8080)
	Status               JobStatus         `json:"status"`
	CreatedAt            time.Time         `json:"createdAt"`
	UpdatedAt            time.Time         `json:"updatedAt"`
	Result               string            `json:"result"` // perhaps some containers will provide a result
	WorkerID             string            `json:"workerId"`
	ErrorMessage         string            `json:"errorMessage"`
	ComputeZone          string            `json:"computeZone"`     // saved as "zone key", we get from Electricity Maps API, e.g "DE" (germany)
	CarbonIntensity      int               `json:"carbonIntensity"` // CO2eq/kWh which are emitted during job execution
	CarbonSaving         int               `json:"carbonSavings"`   // consumption savings compared to the actual consumer location
}
