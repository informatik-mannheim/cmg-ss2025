package ports

import (
	"context"
	"errors"
)

var ErrConsumerNotFound = errors.New("Consumer Gateway not found")

type CreateJobRequest struct {
	ImageID string `json:"image_id"`
	Zone string  `json:"zone"` // is optional
	Param string `json:"params"`
}

type CreateJobResponse struct {
	JobID string `json:"job_id"`
	JobStatus string `json:"status"`
}

type GetJobResult struct {
	JobID string `json:"job_id"`
}

type JobResultResponse struct {
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

type ZoneRequest struct {
	Zone string `json:"zone"`
}
// CLI: "Choose job location [DE, FR, ESP ...]"
type ZoneResponse struct {
	Zone string `json:"zone"`
}



type Api interface {

	CreateJob(req CreateJobRequest, ctx context.Context) (CreateJobResponse, error)
	GetJobResult(jobID string, ctx context.Context) (JobResultResponse, error)

	Login(req ConsumerLoginRequest, ctx context.Context) (LoginResponse, error)
	Register(req ConsumerRegistrationRequest, ctx context.Context) (RegisterResponse, error)

}