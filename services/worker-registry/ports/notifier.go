package ports

import (
	"context"
)

type Notifier interface {
	WorkerChanged(worker Worker, ctx context.Context)
}
