package ports

import (
	"context"
)

type JobService interface {
	UpdateJob(ctx context.Context, req ResultRequest) error
	FetchScheduledJobs(ctx context.Context) ([]Job, error)
}
