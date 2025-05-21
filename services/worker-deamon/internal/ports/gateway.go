package ports

type WorkerGateway interface {
	Register(key string, zone string) (*RegisterResponse, error)
	SendHeartbeat(workerID string, status string) ([]Job, error)
	SendResult(j Job) error
}

type Job struct {
	ID           string `json:"id"`
	Status       string `json:"status"`
	Result       string `json:"result"`
	ErrorMessage string `json:"errorMessage"`
}

type RegisterResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Zone   string `json:"zone"`
	Token  string `json:"token"`
}
