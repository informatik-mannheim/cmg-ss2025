package worker

import (
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/model"
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/ports"
)

type WorkerAdapter struct{}

var _ ports.WorkerAdapter = (*WorkerAdapter)(nil)

func NewWorkerAdapater() *WorkerAdapter {
	return &WorkerAdapter{}
}

func (adapter *WorkerAdapter) AssignWorker(update ports.UpdateWorker) error {
	// FIXME: implement
	panic("unimplemented")
}

func (adpater *WorkerAdapter) GetWorkers() (model.GetWorkersResponse, error) {
	// FIXME: implement
	panic("unimplemented")
}
