package notifier

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/ports"
)

type Notifier struct{}

var _ ports.Notifier = (*Notifier)(nil)

func NewNotifier() *Notifier {
	return &Notifier{}
}

func (n *Notifier) NotifyAssignment(jobId, workerId uuid.UUID) error {
	// TODO: temporary implementation, will probably change, but not in Phase 2
	message := fmt.Sprintf("Job %s and Worker %s assigned successfully\n", jobId.String(), workerId.String())
	log.Print(message)
	return nil
}

func (n *Notifier) NotifyWorkerAssignmentFailed(jobId, workerId uuid.UUID) error {
	// TODO: temporary implementation, will probably change, but not in Phase 2
	message := fmt.Sprintf("Job %s assigned to Worker %s, but Worker assignment failed\n", jobId.String(), workerId.String())
	log.Print(message)
	return nil
}

func (n *Notifier) NotifyAssigmentCorrection(jobId, workerId uuid.UUID) error {
	// TODO: temporary implementation, will probably change, but not in Phase 2
	message := fmt.Sprintf("Corrected failed Worker assignment for Job %s to Worker %s\n", jobId.String(), workerId.String())
	log.Print(message)
	return nil
}
