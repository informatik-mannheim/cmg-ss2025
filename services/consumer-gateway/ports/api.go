package ports

import (
	"context"
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
	Image        ContainerImage    `json:"image"`
	JobName      string            `json:"jobName"`
	CreationZone string            `json:"creationZone"`
	Parameters   map[string]string `json:"parameters"`
	Status       string            `json:"status"`
}

// Returns a singular job
type GetJob struct {
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

type ConsumerLoginRequest struct {
	Secret string `json:"secret"`
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

type Api interface {
	CreateJob(ctx context.Context, req CreateJobRequest) (CreateJobResponse, error)
	GetJobOutcome(ctx context.Context, jobID string) (JobOutcomeResponse, error)
	GetZone(ctx context.Context, req ZoneRequest) (ZoneResponse, error)
	Login(ctx context.Context, req ConsumerLoginRequest) (LoginResponse, error)
}
