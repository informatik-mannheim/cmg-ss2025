package notifier

import (
	"context"

	"github.com/informatik-mannheim/cmg-ss2025/services/worker-registry/ports"
)

type HttpNotifier struct{}

func NewHttpNotifier() *HttpNotifier {
	return &HttpNotifier{}
}

// currently just dummy notifier to make the service runnable
func (n *HttpNotifier) WorkerChanged(worker ports.Worker, ctx context.Context) {
	// TODO implement Notifiermethod
}
