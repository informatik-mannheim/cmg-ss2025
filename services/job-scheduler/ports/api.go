package ports

// api.go - For possible future interaction (like metrics or CLI)

type Api interface {
	ScheduleJob() error
}

type ScheduleResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type ScheduleRequest struct {
	Secret string `json:"secret"`
}
