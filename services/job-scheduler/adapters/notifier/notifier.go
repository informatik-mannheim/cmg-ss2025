package notifier

import (
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/model"
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/ports"
)

type Notifier struct{}

var _ ports.Notifier = (*Notifier)(nil)

func (n *Notifier) NotifyAssignment(job model.Job, worker model.Worker) error {
	// FIXME: implement
	panic("unimplemented")
}

func (n *Notifier) NotifyWorkerAssignmentFailed(job model.Job, worker model.Worker) error {
	// FIXME: implement
	panic("unimplemented")
}

func (n *Notifier) NotifyAssigmentCorrection(job model.Job, worker model.Worker) error {
	// FIXME: implement
	panic("unimplemented")
}
