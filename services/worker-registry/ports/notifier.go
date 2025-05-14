package ports

import (
	"context"
)

type Notifier interface {
	WorkerCreated(worker Worker, ctx context.Context)
	WorkerStatusChanged(worker Worker, ctx context.Context)
}
