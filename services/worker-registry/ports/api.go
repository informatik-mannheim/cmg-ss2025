package ports

import (
	"context"
	"errors"
)

var ErrWorkerNotFound = errors.New("Worker not found")
var ErrUpdatingWorkerFailed = errors.New("invalid status ('AVAILABLE' or 'RUNNING')")
var ErrStoringWorkerFailed = errors.New("storing worker failed due to missing parameters (status or zone)")

type Api interface {
	GetWorkers(status WorkerStatus, zone string, ctx context.Context) ([]Worker, error)
	GetWorkerById(id string, ctx context.Context) (Worker, error)
	CreateWorker(zone string, ctx context.Context) (Worker, error)
	UpdateWorkerStatus(id string, status WorkerStatus, ctx context.Context) (Worker, error)
}
