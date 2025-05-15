package ports

import (
	"context"
)

// JobCreate represents the required fields for job creation
type JobCreate struct {
	JobName      string            `json:"jobName"`
	CreationZone string            `json:"creationZone"`
	Image        ContainerImage    `json:"image"`
	Parameters   map[string]string `json:"parameters"`
}

// SchedulerUpdateData represents data needed for updating a job from the scheduler's perspective
type SchedulerUpdateData struct {
	WorkerID        string    `json:"workerId"`
	ComputeZone     string    `json:"computeZone"`
	CarbonIntensity int       `json:"carbonIntensity"`
	CarbonSaving    int       `json:"carbonSavings"`
	Status          JobStatus `json:"status"`
}

// WorkerDaemonUpdateData represents data needed for updating a job from the worker daemon's perspective
type WorkerDaemonUpdateData struct {
	Status       JobStatus `json:"status"`
	Result       string    `json:"result"`
	ErrorMessage string    `json:"errorMessage"`
}

// JobOutcome represents the outcome of a job
type JobOutcome struct {
	JobName         string    `json:"jobName"`
	Status          JobStatus `json:"status"`
	Result          string    `json:"result"`
	ErrorMessage    string    `json:"errorMessage"`
	ComputeZone     string    `json:"computeZone"`
	CarbonIntensity int       `json:"carbonIntensity"`
	CarbonSavings   int       `json:"carbonSavings"`
}

// JobService defines interfaces for interacting with Job resources
type JobService interface {
	// GetJobs retrieves jobs, optionally filtered by status
	GetJobs(ctx context.Context, status []JobStatus) ([]Job, error)

	// CreateJob creates a new job in the queue
	CreateJob(ctx context.Context, job JobCreate) (Job, error)

	// GetJob retrieves a specific job by its ID
	GetJob(ctx context.Context, id string) (Job, error)

	// GetJobOutcome retrieves detailed result and metadata of a job by its ID
	GetJobOutcome(ctx context.Context, id string) (JobOutcome, error)

	// UpdateJobScheduler updates job properties from a scheduler's perspective
	UpdateJobScheduler(ctx context.Context, id string, data SchedulerUpdateData) (Job, error)

	// UpdateJobWorkerDaemon updates job properties from a worker daemon's perspective
	UpdateJobWorkerDaemon(ctx context.Context, id string, data WorkerDaemonUpdateData) (Job, error)
}
