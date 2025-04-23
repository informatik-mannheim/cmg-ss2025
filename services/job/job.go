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

// Json tags helps by (de-)serializing, json:"id" -> "id":"1234", functionality imported by "encoding/json" package
// If fields value ist empty/nil and its tagged with omitepty, json ignores it

type Job struct {
	ID                   string         `json:"id"`
	UserID               string         `json:"userId"`
	JobName              string         `json:"jobName"`
	Payload              string         `json:"payload"`
	AdjustmentParameters []string       `json:"parameters"`
	Priority             JobPriority    `json:"priority"`
	Status               JobStatus      `json:"status"`
	WorkerID             string         `json:"workerId,omitempty"`
	Constraints          JobConstraints `json:"constraints"`
	Result               *string        `json:"result,omitempty"` // e.g buffered result
	ErrorMessage         *string        `json:"errorMessage,omitempty"`
	RetryCount           int            `json:"retryCount"`
	MaxRetries           int            `json:"maxRetries"`
	CreatedAt            time.Time      `json:"createdAt"`
	StartedAt            *time.Time     `json:"startedAt,omitempty"`
	CompletedAt          *time.Time     `json:"completedAt,omitempty"`
	LastUpdatedAt        time.Time      `json:"lastUpdatedAt"`
	UsedCarbonItensity   int            `json:"UsedCarbonItensity`
	// AvailableJobs		 bool ja/nein?
	// Filterung mithilfe von Query Parametern?
}

type JobConstraints struct {
	MaxRuntime       int      `json:"maxRuntime"`       // e.g in seconds
	PreferredRegions []string `json:"preferredRegions"` // city, continent?
}
