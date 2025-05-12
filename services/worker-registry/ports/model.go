package ports

type WorkerStatus string

const (
	StatusAvailable WorkerStatus = "AVAILABLE" // default value for new worker
	StatusRunning   WorkerStatus = "RUNNING"   // set by Job Scheduler
)

type UpdateWorkerStatusRequest struct {
	Status WorkerStatus `json:"status"`
}

type Worker struct {
	Id     string       `json:"id"`
	Status WorkerStatus `json:"status"`
	Zone   string       `json:"zone"`
}

type Zone struct {
	Code string `json:"code"`
	Name string `json:"name"`
}
