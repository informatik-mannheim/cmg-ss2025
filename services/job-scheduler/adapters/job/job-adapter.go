package job

import (
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/model"
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/ports"
)

type JobAdapter struct{}

var _ ports.JobAdapter = (*JobAdapter)(nil)

func NewJobAdapter() *JobAdapter {
	return &JobAdapter{}
}

func (adapter *JobAdapter) GetJobs() (model.GetJobsResponse, error) {
	// FIXME: implement
	panic("unimplemented")
}

func (adapter *JobAdapter) AssignJob(update ports.UpdateJob) error {
	// FIXME: implement
	panic("unimplemented")
}
