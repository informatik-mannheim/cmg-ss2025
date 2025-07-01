package cli

type CreateJobRequest struct {
	JobName      string            `json:"jobName"`
	CreationZone string            `json:"creationZone"`
	Image        ContainerImage    `json:"image"`
	Parameters   map[string]string `json:"parameters"`
}
type ContainerImage struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type TokenResponse struct {
	Token string `json:"secret"`
}
