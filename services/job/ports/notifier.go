package ports

import (
	"context"
)

type Notifier interface {
	JobChanged(job Job, ctx context.Context)
}
