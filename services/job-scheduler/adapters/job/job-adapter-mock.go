package job

import (
	"fmt"

	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/model"
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/ports"
)

type JobAdapterMock struct {
	shouldGetJobsFail   bool
	shouldGetJobsEmpty  bool
	shouldAssingJobFail bool
}

var _ ports.JobAdapter = (*JobAdapterMock)(nil)

func NewJobAdapterMock(shouldGetJobsFail, shouldGetJobsEmpty, shouldAssingJobFail bool) *JobAdapterMock {
	return &JobAdapterMock{
		shouldGetJobsFail:   shouldGetJobsFail,
		shouldAssingJobFail: shouldAssingJobFail,
		shouldGetJobsEmpty:  shouldGetJobsEmpty,
	}
}

func (adapter *JobAdapterMock) GetJobs() (model.GetJobsResponse, error) {
	if adapter.shouldGetJobsFail {
		return nil, fmt.Errorf("some job get error")
	}
	if adapter.shouldGetJobsEmpty {
		return model.GetJobsResponse{}, nil
	}
	return MockJobs, nil
}

func (adapter *JobAdapterMock) AssignJob(update ports.UpdateJob) error {
	if adapter.shouldAssingJobFail {
		return fmt.Errorf("some job assignment error")
	}
	return nil
}
