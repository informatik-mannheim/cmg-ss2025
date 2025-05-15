package ports

import (
	"context"
)

type RegistryService interface {
	RegisterWorker(ctx context.Context, req RegisterRequest) error
	UpdateWorkerStatus(ctx context.Context, req HeartbeatRequest) error
}
