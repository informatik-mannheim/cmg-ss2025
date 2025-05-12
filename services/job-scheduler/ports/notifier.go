package ports

import "github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/model"

type Notifier interface {
	NotifyAssignment(job model.Job, worker model.Worker) error             // for when both job and worker got assigned
	NotifyWorkerAssignmentFailed(job model.Job, worker model.Worker) error // for when the job got assigned, but the worker for whatever reason not
	NotifyAssigmentCorrection(job model.Job, worker model.Worker) error    // for when an worker gets assigned after it failed in a previous cycle
}
