package ports

import (
	"context"
	"fmt"
)

func NewErrWorkerNotFound(id string) error {
	return fmt.Errorf("Worker with ID %v not found", id)
}

func NewErrUpdatingWorkerFailed(id string) error {
	return fmt.Errorf("invalid status ('AVAILABLE' or 'RUNNING') for worker with ID %v", id)
}

func NewErrCreatingWorkerFailed() error {
	return fmt.Errorf("creating worker failed due to missing parameter 'zone'")
}

type Api interface {
	GetWorkers(status WorkerStatus, zone string, ctx context.Context) ([]Worker, error)
	GetWorkerById(id string, ctx context.Context) (Worker, error)
	CreateWorker(zone string, ctx context.Context) (Worker, error)
	UpdateWorkerStatus(id string, status WorkerStatus, ctx context.Context) (Worker, error)
}
