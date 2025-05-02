package core_test

import (
	"context"

	"github.com/informatik-mannheim/cmg-ss2025/services/worker-registry/ports"
)

type MockRepo struct {
	workers      []ports.Worker
	storedWorker ports.Worker
	err          error
}

type MockNotifier struct {
	worker    ports.Worker
	callcount int
}

func (notifier *MockNotifier) WorkerChanged(worker ports.Worker, ctx context.Context) {
	notifier.worker = worker
	notifier.callcount++
}

func (repo *MockRepo) GetWorkers(status, zone string, ctx context.Context) ([]ports.Worker, error) {
	if repo.err != nil {
		return nil, repo.err
	}
	return repo.workers, nil
}

func (repo *MockRepo) GetWorkerById(id string, ctx context.Context) (ports.Worker, error) {
	if repo.err != nil {
		return ports.Worker{}, repo.err
	}
	return repo.storedWorker, nil
}

func (repo *MockRepo) StoreWorker(worker ports.Worker, ctx context.Context) error {
	if repo.err != nil {
		return repo.err
	}
	repo.storedWorker = worker
	return nil
}

func (repo *MockRepo) UpdateWorkerStatus(id, status string, ctx context.Context) (ports.Worker, error) {
	if repo.err != nil {
		return ports.Worker{}, repo.err
	}

	return ports.Worker{Id: id, Status: status}, nil
}

var _ ports.Repo = (*MockRepo)(nil)
var _ ports.Notifier = (*MockNotifier)(nil)
