package repo_in_memory

import (
	"context"

	"github.com/informatik-mannheim/cmg-ss2025/services/worker-registry/ports"
)

type Repo struct {
	workers map[string]ports.Worker
}

var _ ports.Repo = (*Repo)(nil)

func NewRepo() *Repo {
	return &Repo{
		workers: make(map[string]ports.Worker),
	}
}

func (r *Repo) GetWorkers(status, zone string, ctx context.Context) ([]ports.Worker, error) {
	var workers []ports.Worker
	for _, worker := range r.workers {
		if (status == "" || worker.Status == status) && (zone == "" || worker.Zone == zone) {
			workers = append(workers, worker)
		}
	}
	return workers, nil
}

func (r *Repo) GetWorkerById(id string, ctx context.Context) (ports.Worker, error) {
	worker, ok := r.workers[id]
	if !ok {
		return ports.Worker{}, ports.ErrWorkerNotFound
	}
	return worker, nil
}

func (r *Repo) StoreWorker(worker ports.Worker, ctx context.Context) error {
	r.workers[worker.Id] = worker
	return nil
}

func (r *Repo) UpdateWorkerStatus(id, status string, ctx context.Context) (ports.Worker, error) {
	worker, ok := r.workers[id]
	if !ok {
		return ports.Worker{}, ports.ErrWorkerNotFound
	}

	worker.Status = status
	r.workers[id] = worker
	return worker, nil
}
