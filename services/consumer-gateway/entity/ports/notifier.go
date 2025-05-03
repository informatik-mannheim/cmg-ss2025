package ports

import (
	"context"
)

type Notifier interface {
	ConsumerChanged(consumer Consumer, ctx context.Context)
}
