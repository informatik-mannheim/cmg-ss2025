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
	ImageID      string            `json:"image_id"`
	JobName      string            `json:"jobName"`
	CreationZone string            `json:"creationZone"`
	Parameters   map[string]string `json:"parameters"`
}

type CreateJobResponse struct {
	ImageID      string            `json:"image_id"`
	JobName      string            `json:"jobName"`
	CreationZone string            `json:"creationZone"`
	Parameters   map[string]string `json:"parameters"`
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
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Secret string `json:"secret"`
}

type ConsumerRegistrationRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type RegisterResponse struct {
	Secret string `json:"secret"`
}

// Get all available zones
type ZoneRequest struct {
	Zone string `json:"zone"`
}

type ZoneResponse struct {
	Zone string `json:"zone"`
}

type ZoneClient interface {
	GetZone(ctx context.Context, req ZoneRequest) (ZoneResponse, error)
}

type LoginClient interface {
	Login(ctx context.Context, req ConsumerLoginRequest) (LoginResponse, error)
}

type JobClient interface {
	GetJobOutcome(ctx context.Context, jobID string) (JobOutcomeResponse, error)
	CreateJob(ctx context.Context, req CreateJobRequest) (CreateJobResponse, error)
}
