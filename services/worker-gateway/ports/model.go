package ports

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

// Job represents a compute job from Job-Service
type Job struct {
	ID       string
	WorkerID string
	Status   string
	Result   string
	ErrorMsg string
}
