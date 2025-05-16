package ports

import (
	"context"
)

type Api interface {
	Heartbeat(ctx context.Context, req HeartbeatRequest) ([]Job, error)
	Result(ctx context.Context, result ResultRequest) error
	Register(ctx context.Context, req RegisterRequest) error
}

// incoming heartbeat from a worker
type HeartbeatRequest struct {
	WorkerID string `json:"workerId"`
	Status   string `json:"status"` // AVAILABLE or COMPUTING
}

// new worker registration
type RegisterRequest struct {
	ID       string `json:"id"`
	Key      string `json:"key"`
	Location string `json:"location"`
}

// a finished job result
type ResultRequest struct {
	JobID        string `json:"jobId"`
	Status       string `json:"status"`
	Result       string `json:"result"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}
