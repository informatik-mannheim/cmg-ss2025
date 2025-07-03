package ports

import (
	"context"
)

type JobService interface {
	UpdateJob(ctx context.Context, req ResultRequest, token string) error
	FetchScheduledJobs(ctx context.Context, token string) ([]Job, error)
}

type Job struct {
	ID                   string            `json:"id"`
	WorkerID             string            `json:"workerId"`
	Image                ContainerImage    `json:"image"`
	AdjustmentParameters map[string]string `json:"adjustmentParameters"`
	Status               string            `json:"status"`
	Result               string            `json:"result"`
	ErrorMessage         string            `json:"errorMessage"`
}

type ContainerImage struct {
	Name    string `json:"name" db:"image_name"`
	Version string `json:"version" db:"image_version"`
}
