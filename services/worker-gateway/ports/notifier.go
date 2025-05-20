package ports

import (
	"context"
)

type Notifier interface {
	RegisterWorker(ctx context.Context, req RegisterRequest) error
	UpdateJob(ctx context.Context, req ResultRequest) error
	UpdateWorkerStatus(ctx context.Context, req HeartbeatRequest) error
	FetchScheduledJobs(ctx context.Context) ([]Job, error)
}
