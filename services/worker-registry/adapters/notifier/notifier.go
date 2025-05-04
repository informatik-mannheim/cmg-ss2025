package notifier

import (
	"context"

	"github.com/informatik-mannheim/cmg-ss2025/services/worker-registry/ports"
)

type Notifier interface {
	WorkerChanged(worker ports.Worker, ctx context.Context)
}

type HttpNotifier struct{}

func NewHttpNotifier() *HttpNotifier {
	return &HttpNotifier{}
}

var _ Notifier = (*HttpNotifier)(nil)

// currently just dummy notifier to make the service runnable
func (n *HttpNotifier) WorkerChanged(worker ports.Worker, ctx context.Context) {
	// TODO implement Notifier method
}
