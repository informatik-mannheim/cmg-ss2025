package cli

type JobStatus string

type CreateJobRequest struct {
	JobName      string            `json:"jobName"`
	CreationZone string            `json:"creationZone"`
	Image        ContainerImage    `json:"image"`
	Parameters   map[string]string `json:"parameters"`
}

type CreateJobResponse struct {
	Image        ContainerImage    `json:"image"`
	JobName      string            `json:"jobName"`
	CreationZone string            `json:"creationZone"`
	Parameters   map[string]string `json:"parameters"`
	Status       string            `json:"status"`
}

type JobOutcomeResponse struct {
	JobName         string    `json:"jobName"`
	Status          JobStatus `json:"status"`
	Result          string    `json:"result"`
	ErrorMessage    string    `json:"errorMessage"`
	ComputeZone     string    `json:"computeZone"`
	CarbonIntensity int       `json:"carbonIntensity"`
	CarbonSavings   int       `json:"carbonSavings"`
}

type ContainerImage struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type TokenResponse struct {
	Token string `json:"secret"`
}
