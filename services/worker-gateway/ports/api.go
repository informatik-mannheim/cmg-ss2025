package ports

import (
	"context"
)

type Api interface {
	Heartbeat(ctx context.Context, req HeartbeatRequest, token string) ([]Job, error)
	Result(ctx context.Context, result ResultRequest, token string) error
	Register(ctx context.Context, req RegisterRequest) (*RegisterRespose, error)
}

// incoming heartbeat from a worker
type HeartbeatRequest struct {
	WorkerID string `json:"workerId"`
	Status   string `json:"status"` // AVAILABLE or RUNNING
}

// new worker registration
type RegisterRequest struct {
	Key  string `json:"key"`
	Zone string `json:"zone"`
}

// a finished job result
type ResultRequest struct {
	JobID        string `json:"jobId"`
	Status       string `json:"status"`
	Result       string `json:"result"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}

type RegisterRespose struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Zone   string `json:"zone"`
	Token  string `json:"token"`
}
