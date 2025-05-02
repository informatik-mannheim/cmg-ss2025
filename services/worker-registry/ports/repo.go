package ports

import (
	"context"
)

type Repo interface {
	GetWorkers(status, zone string, ctx context.Context) ([]Worker, error)
	GetWorkerById(id string, ctx context.Context) (Worker, error)
	StoreWorker(worker Worker, ctx context.Context) error
	UpdateWorkerStatus(id, status string, ctx context.Context) (Worker, error)
}
