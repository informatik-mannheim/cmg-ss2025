package ports

import (
	"context"
)

type Notifier interface {
	EntityChanged(entity Entity, ctx context.Context)
}
