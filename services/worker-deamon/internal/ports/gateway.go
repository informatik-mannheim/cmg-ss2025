package ports

type WorkerGateway interface {
	Register(key string, zone string) (*RegisterResponse, error)
	SendHeartbeat(workerID string, status string, token string) ([]Job, error)
	SendResult(j Job, token string) error
}

type Job struct {
	ID                   string            `json:"id"`
	Image                ContainerImage    `json:"image"`
	AdjustmentParameters map[string]string `json:"adjustmentParameters"`
	Status               string            `json:"status"`
	Result               string            `json:"result"`
	ErrorMessage         string            `json:"errorMessage"`
}

type RegisterResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Zone   string `json:"zone"`
	Token  string `json:"token"`
}

type ContainerImage struct {
	Name    string `json:"name" db:"image_name"`
	Version string `json:"version" db:"image_version"`
}
