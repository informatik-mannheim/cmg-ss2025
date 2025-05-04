package ports

type WorkerStatus string

const (
	StatusAvailable WorkerStatus = "AVAILABLE" // default value for new worker
	StatusRunning   WorkerStatus = "RUNNING"   // set by Job Scheduler
)

type Worker struct {
	Id     string       `json:"id"`
	Status WorkerStatus `json:"status"`
	Zone   string       `json:"zone"`
}
