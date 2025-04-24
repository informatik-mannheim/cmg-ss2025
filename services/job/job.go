package models

import "time"

type JobStatus string

const (
	StatusQueued    JobStatus = "queued"
	StatusRunning   JobStatus = "running"
	StatusCompleted JobStatus = "completed"
	StatusFailed    JobStatus = "failed"
	StatusCancelled JobStatus = "cancelled"
)

type JobPriority int

// Easier to compare later on
const (
	Low JobPriority = iota
	Middle
	High
)

type JobConstraints struct {
	MaxRuntime        int      `json:"maxRuntime"`        // e.g in seconds
	PreferredLocation []string `json:"preferredLocation"` // city, continent?
}

type ContainerImage struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Version string `json:"version"`
}

// Json tags helps by (de-)serializing, json:"id" -> "id":"1234", functionality imported by "encoding/json" package
// If fields value ist empty/nil (needs pointer) and its tagged with omitempty, json ignores it

type Job struct {
	JobID                string         `json:"jobId"`
	JobName              string         `json:"jobName"`
	UserID               string         `json:"userId"`
	Image                ContainerImage `json:"image"`
	AdjustmentParameters []string       `json:"parameters"`
	Priority             JobPriority    `json:"priority"` // 0,1,2
	Status               JobStatus      `json:"status"`
	CreatedAt            time.Time      `json:"createdAt"`
	LastUpdatedAt        time.Time      `json:"lastUpdatedAt"`
	WorkerID             *string        `json:"workerId,omitempty"`
	ErrorMessage         *string        `json:"errorMessage,omitempty"`
	ComputeLocation      *string        `json:"computeLocation,omitempty"`
	CarbonIntensity      *int           `json:"carbonIntensity,omitempty"` // grams CO2 per kWH

	// Optional fields:
	Constraints *JobConstraints `json:"constraints,omitempty"`
	Result      *string         `json:"result,omitempty"` // perhaps some containers will provide a result
	MaxRetries  *int            `json:"maxRetries,omitempty"`
	StartedAt   *time.Time      `json:"startedAt,omitempty"`
	CompletedAt *time.Time      `json:"completedAt,omitempty"`
}
