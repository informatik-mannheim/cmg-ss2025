package ports

import (
	"context"
	"errors"
)

var ErrWorkerNotFound = errors.New("Worker not found")

type Api interface {
	GetWorkers(status, zone string, ctx context.Context) ([]Worker, error)
	GetWorkerById(id string, ctx context.Context) (Worker, error)
	CreateWorker(zone string, ctx context.Context) (Worker, error)
	UpdateWorkerStatus(id, status string, ctx context.Context) (Worker, error)
}
