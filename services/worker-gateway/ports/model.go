package ports

// Job represents a compute job from Job-Service
type Job struct {
	ID       string
	WorkerID string
	Status   string
	Result   string
	ErrorMsg string
}
