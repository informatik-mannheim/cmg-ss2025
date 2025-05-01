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

type ConsumerLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`

}

type LoginResponse struct {
	Token string `json:"token"`
}

type ConsumerRegistrationRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role string `json:"role"`
}

type RegisterResponse struct {
	Token string `json:"token"`
}

// Request to /me
type MeRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type MeResponse struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role string `json:"role"`
}

type Api interface {

	CreateJobRequest(req CreateJobRequest, ctx context.Context) (CreateJobResponse, error)
	GetJobStatus(jobID string, ctx context.Context) (CreateJobRequest, error)

	ConsumerLoginRequest(req ConsumerLoginRequest, ctx context.Context) (LoginResponse, error)
	ConsumerRegisterRequest(req ConsumerRegistrationRequest, ctx context.Context) (RegisterResponse, error)
	MeRequest(ctx context.Context) (MeResponse, error)
}