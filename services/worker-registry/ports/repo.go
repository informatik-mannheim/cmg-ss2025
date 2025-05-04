package ports

import (
	"context"
)

type Repo interface {
	GetWorkers(status WorkerStatus, zone string, ctx context.Context) ([]Worker, error)
	GetWorkerById(id string, ctx context.Context) (Worker, error)
	CreateWorker(worker Worker, ctx context.Context) error
	UpdateWorkerStatus(id string, status WorkerStatus, ctx context.Context) (Worker, error)
}
