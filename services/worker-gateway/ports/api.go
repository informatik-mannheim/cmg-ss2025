package ports

import (
	"context"
)

type Api interface {
	Heartbeat(ctx context.Context, req HeartbeatRequest) ([]Job, error)
	Result(ctx context.Context, result ResultRequest) error
	Register(ctx context.Context, req RegisterRequest) error
}
