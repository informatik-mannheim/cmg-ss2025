package ports

import (
	"context"
)

type JobStorage interface {
	GetJobs(ctx context.Context, status []JobStatus) ([]Job, error)
	CreateJob(ctx context.Context, job Job) error
	GetJob(ctx context.Context, id string) (Job, error)
	UpdateJob(ctx context.Context, id string, job Job) (Job, error)
}
