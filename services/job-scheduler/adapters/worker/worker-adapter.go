package worker

import (
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/model"
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/ports"
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/utils"
)

type WorkerAdapter struct {
	environments model.Environments
}

var _ ports.WorkerAdapter = (*WorkerAdapter)(nil)

func NewWorkerAdapter(environments model.Environments) *WorkerAdapter {
	return &WorkerAdapter{
		environments: environments,
	}
}

func (adapter *WorkerAdapter) GetWorkers() (model.GetWorkersResponse, error) {
	// For now its kept simple and return an error as soon as it gets one, changes in Phase 3
	endpoint := model.GetWorkersEndpoint(adapter.environments.WorkerRegestryUrl)

	data, err := utils.GetRequest[model.GetWorkersResponse](endpoint)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (adapter *WorkerAdapter) AssignWorker(update ports.UpdateWorker) error {
	// For now its kept simple and return an error as soon as it gets one, changes in Phase 3
	endpoint := model.PutWorkerStatusEndpoint(adapter.environments.WorkerRegestryUrl, update.ID)

	payload := model.UpdateWorkerPayload{
		WorkerStatus: model.WorkerStatusRunning,
	}

	_, err := utils.PutRequest[model.UpdateWorkerPayload, model.UpdateWorkerResponse](endpoint, payload)
	if err != nil {
		return err
	}

	return nil
}
