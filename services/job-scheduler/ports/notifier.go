package ports

import (
	"github.com/google/uuid"
)

type Notifier interface {
	NotifyAssignment(jobId, workerId uuid.UUID) error             // for when both job and worker got assigned
	NotifyWorkerAssignmentFailed(jobId, workerId uuid.UUID) error // for when the job got assigned, but the worker for whatever reason not
	NotifyAssigmentCorrection(jobId, workerId uuid.UUID) error    // for when an worker gets assigned after it failed in a previous cycle
}
