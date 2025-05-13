package worker

import (
	"fmt"

	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/model"
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/ports"
)

type WorkerAdapterMock struct {
	shouldGetWorkersFail    bool
	shouldGetWorkersEmpty   bool
	shouldAssignWorkersFail bool
}

var _ ports.WorkerAdapter = (*WorkerAdapterMock)(nil)

func NewWorkerAdapterMock(shouldGetWorkersFail, shouldGetWorkersEmpty, shouldAssignWorkersFail bool) *WorkerAdapterMock {
	return &WorkerAdapterMock{
		shouldGetWorkersFail:    shouldGetWorkersFail,
		shouldAssignWorkersFail: shouldAssignWorkersFail,
		shouldGetWorkersEmpty:   shouldGetWorkersEmpty,
	}
}

func (adapter *WorkerAdapterMock) GetWorkers() (model.GetWorkersResponse, error) {
	if adapter.shouldGetWorkersFail {
		return nil, fmt.Errorf("some worker get error")
	}
	if adapter.shouldGetWorkersEmpty {
		return model.GetWorkersResponse{}, nil
	}
	return MockWorkers, nil
}

func (adapter *WorkerAdapterMock) AssignWorker(update ports.UpdateWorker) error {
	if adapter.shouldAssignWorkersFail {
		return fmt.Errorf("some worker assignment error")
	}

	return nil
}
