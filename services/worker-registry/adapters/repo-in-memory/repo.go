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
	matchingWorkers := []ports.Worker{}
	for _, worker := range r.workers {
		if (status == "" || worker.Status == status) && (zone == "" || worker.Zone == zone) {
			matchingWorkers = append(matchingWorkers, worker)
		}
	}
	return matchingWorkers, nil
}

func (r *Repo) GetWorkerById(id string, ctx context.Context) (ports.Worker, error) {
	worker, ok := r.workers[id]
	if !ok {
		return ports.Worker{}, ports.ErrWorkerNotFound
	}
	return worker, nil
}

func (r *Repo) StoreWorker(worker ports.Worker, ctx context.Context) error {
	if worker.Status == "" || worker.Zone == "" {
		return ports.ErrStoringWorkerFailed
	}
	r.workers[worker.Id] = worker
	return nil
}

func (r *Repo) UpdateWorkerStatus(id, status string, ctx context.Context) (ports.Worker, error) {
	worker, ok := r.workers[id]
	if !ok {
		return ports.Worker{}, ports.ErrWorkerNotFound
	}
	if !isValidStatus(status) {
		return ports.Worker{}, ports.ErrUpdatingWorkerFailed
	}

	worker.Status = status
	r.workers[id] = worker
	return worker, nil
}

func isValidStatus(status string) bool {
	return status == "AVAILABLE" || status == "RUNNING"
}
