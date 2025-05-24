package ports

import (
	"context"
)

type RegistryService interface {
	RegisterWorker(ctx context.Context, req RegisterRequest) (*RegisterRespose, error)
	UpdateWorkerStatus(ctx context.Context, req HeartbeatRequest) error
}
