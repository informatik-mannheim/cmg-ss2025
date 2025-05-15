package notifier

import (
	"context"
	"log"

	"github.com/informatik-mannheim/cmg-ss2025/services/worker-registry/ports"
)

type Notifier struct{}

func NewNotifier() ports.Notifier {
	return &Notifier{}
}

func (n *Notifier) WorkerCreated(worker ports.Worker, ctx context.Context) {
	log.Printf("[Notifier] New worker created: ID=%s, STATUS=%s, ZONE=%s ", worker.Id, worker.Status, worker.Zone)
}

func (n *Notifier) WorkerStatusChanged(worker ports.Worker, ctx context.Context) {
	log.Printf("[Notifier] Changed status from Worker with ID '%s' to status '%s'.", worker.Id, worker.Status)
}
