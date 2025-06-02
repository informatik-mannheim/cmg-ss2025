package ports

import (
	"errors"
)

var ErrNotFound = errors.New("not found")
var ErrUnauthorized = errors.New("unauthorized")
var ErrBadRequest = errors.New("bad Request")
var ErrInvalidInput = errors.New("invalid Input")

type CreateJobRequest struct {
	JobName      string            `json:"jobName"`
	CreationZone string            `json:"creationZone"`
	ImageID      ContainerImage    `json:"image"`
	Parameters   map[string]string `json:"parameters"`
}

type CreateJobResponse struct {
	Image        string            `json:"image_id"`
	JobName      string            `json:"job_name"`
	CreationZone string            `json:"creation_zone"`
	Parameters   map[string]string `json:"parameters"`
	Status       string            `json:"status"`
}

// Returns a singular job
type GetJob struct {
}

type JobOutcomeResponse struct {
	JobName         string    `json:"job_name"`
	Status          JobStatus `json:"status"`
	Result          string    `json:"result"`
	ErrorMessage    string    `json:"error_message"`
	ComputeZone     string    `json:"compute_zone"`
	CarbonIntensity int       `json:"carbon_intensity"`
	CarbonSavings   int       `json:"carbon_savings"`
}

type ConsumerLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Secret string `json:"secret"`
}

// Get all available zones
type ZoneRequest struct {
	Zone string `json:"zone"`
}

type ZoneResponse struct {
	Zone string `json:"zone"`
}
