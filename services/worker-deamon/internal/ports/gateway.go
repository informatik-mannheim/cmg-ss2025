package ports

type GatewayClient interface {
	Register(workerID, key, location string) error
	SendHeartbeat(workerID string, status string) ([]Job, error)
	SendResult(j Job) error
}

type Job struct {
	ID           string `json:"id"`
	Status       string `json:"status"`
	Result       string `json:"result"`
	ErrorMessage string `json:"errorMessage"`
}
