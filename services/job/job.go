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

// Json tags helps by (de-)serializing, json:"id" -> "id":"1234", functionality imported by "encoding/json" package
// If fields value ist empty/nil (needs pointer) and its tagged with omitempty, json ignores it

type Job struct {
	JobID                 string      `json:"jobId"`
	JobName               string      `json:"jobName"`
	UserID                string      `json:"userId"`
	ImageID               string      `json:"imageId"`
	ImageName             string      `json:"imageName"`
	ImageVersion          string      `json:"imageVersion"` // e.g :latest
	AdjustmentParameters  []string    `json:"parameters"`
	Priority              JobPriority `json:"priority"` // 0,1,2
	Status                JobStatus   `json:"status"`
	CreatedAt             time.Time   `json:"createdAt"`
	WorkerID              *string     `json:"workerId,omitempty"`
	ErrorMessage          *string     `json:"errorMessage,omitempty"`
	ComputeLocation       *string     `json:"computeLocation,omitempty"`
	CarbonItensity        *int        `json:"carbonItensity,omitempty"` // grams CO2 per kWH
	JobsAvailable         bool        `json:"jobsavailable"`
	NumberOfJobsAvailable int         `json:"numberOfJobsAvailable"`

	// Optional fields:

	Constraints   *JobConstraints `json:"constraints,omitempty"`
	Result        *string         `json:"result,omitempty"` // perhaps some containers will provide a result
	MaxRetries    int             `json:"maxRetries"`
	StartedAt     *time.Time      `json:"startedAt,omitempty"`
	CompletedAt   *time.Time      `json:"completedAt,omitempty"`
	LastUpdatedAt time.Time       `json:"lastUpdatedAt"`
}
