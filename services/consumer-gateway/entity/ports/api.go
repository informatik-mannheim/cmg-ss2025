package ports

import (
	"context"
	"errors"
)

var ErrConsumerNotFound = errors.New("Consumer Gateway not found")

type CreateJobRequest struct {
	ImageID string `json:"image_id"`
	Location string  `json:"location"`
}

type CreateJobResponse struct {
	JobID string `json:"job_id"`
	JobStatus string `json:"status"`
}
type GetJobStatus struct {
	JobID string `json:"job_id"`

}

type JobStatusResponse struct {
	JobID string `json:"job_id"`
	JobStatus string `json:"status"`

}

type ConsumerLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`

}

type LoginResponse struct {
	Token string `json:"token"`
}

type ConsumerRegistration struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role string `json:"role"`
}

type RegisterResponse struct {
	Token string `json:"token"`
}

// Request to /me
type Me struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type MeResponse struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role string `json:"role"`
}

type Api interface {

	CreateJob(req CreateJobRequest, ctx context.Context) (CreateJobResponse, error)
	GetJobStatus(jobID string, ctx context.Context) (CreateJobRequest, error)

	Login(req ConsumerLogin, ctx context.Context) (LoginResponse, error)
	Register(req ConsumerRegistration, ctx context.Context) (RegisterResponse, error)
	GetCurrentUser(ctx context.Context) (MeResponse, error)
}