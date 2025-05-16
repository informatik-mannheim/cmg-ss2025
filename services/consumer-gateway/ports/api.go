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
	ImageID string `json:"image_id"`
	Zone    string `json:"zone"` // is optional
	Param   string `json:"params"`
}

type CreateJobResponse struct {
	ImageID   string `json:"image_id`
	Zone      string `json:"zone"`
	Param     string `json:"params"`
	JobStatus string `json:"job_status"`
}

type GetJobResult struct {
	ImageID string `json:"image_id"`
}

type JobResultResponse struct {
	ImageID   string `json:"image_id"`
	JobStatus string `json:"status"`
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

type Api interface {
	CreateJob(req CreateJobRequest, ctx context.Context) (CreateJobResponse, error)
	GetJobResult(ImageID string, ctx context.Context) (JobResultResponse, error)

	// Get available zones from carbon intesity provider
	GetZone(req ZoneRequest, ctx context.Context) (ZoneResponse, error)

	Login(req ConsumerLoginRequest, ctx context.Context) (LoginResponse, error)
	Register(req ConsumerRegistrationRequest, ctx context.Context) (RegisterResponse, error)
}
