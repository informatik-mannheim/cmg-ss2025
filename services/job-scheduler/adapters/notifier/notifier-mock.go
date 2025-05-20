package notifier

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/ports"
)

type NotifierMock struct {
	shouldNotifyAssignmentFail       bool
	shouldNotifyWorkerAssignmentFail bool
	shouldNotifyAssignmentCorrection bool
}

var _ ports.Notifier = (*NotifierMock)(nil)

func NewNotifierMock(
	shouldNotifyAssignmentFail bool,
	shouldNotifyWorkerAssignmentFail bool,
	shouldNotifyAssigmentCorrection bool,
) *NotifierMock {
	return &NotifierMock{
		shouldNotifyAssignmentFail:       shouldNotifyAssignmentFail,
		shouldNotifyWorkerAssignmentFail: shouldNotifyWorkerAssignmentFail,
		shouldNotifyAssignmentCorrection: shouldNotifyAssigmentCorrection,
	}
}

func (n *NotifierMock) NotifyAssignment(jobId, workerId uuid.UUID) error {
	if n.shouldNotifyAssignmentFail {
		return fmt.Errorf("NotifyAssignment failed")
	}
	return nil
}

func (n *NotifierMock) NotifyWorkerAssignmentFailed(jobId, workerId uuid.UUID) error {
	if n.shouldNotifyWorkerAssignmentFail {
		return fmt.Errorf("NotifyWorkerAssignmentFailed failed")
	}
	return nil
}

func (n *NotifierMock) NotifyAssignmentCorrection(jobId, workerId uuid.UUID) error {
	if n.shouldNotifyAssignmentCorrection {
		return fmt.Errorf("NotifyAssigmentCorrection failed")
	}
	return nil
}
