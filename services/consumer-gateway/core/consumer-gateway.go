package core

import (
	"context"

	"github.com/informatik-mannheim/cmg-ss2025/services/consumer-gateway/ports"
)

type ConsumerService struct {
}

func NewConsumerService() *ConsumerService {
	return &ConsumerService{}
}

func (s *ConsumerService) CreateJob(req ports.CreateJobRequest, ctx context.Context) (ports.CreateJobResponse, error) {
	if req.ImageID == "" || req.Zone == "" || req.Param == "invalid" {
		return ports.CreateJobResponse{}, ports.ErrInvalidInput
	}
	return ports.CreateJobResponse{
		ImageID: req.ImageID,
		Zone:    req.Zone,
		Param:   req.Param,
<<<<<<< HEAD
		Status: req.JobStatus
=======
		JobStatus: "queued",
>>>>>>> origin/main
	}, nil
}

func (s *ConsumerService) GetJobResult(jobID string, ctx context.Context) (ports.JobResultResponse, error) {
	user, ok := ctx.Value("user").(string)
	if !ok || user != "alice" {
		return ports.JobResultResponse{}, ports.ErrNotFound
	}
	return ports.JobResultResponse{
		ImageID:   jobID,
		JobStatus: "completed",
	}, nil
}

func (s *ConsumerService) GetZone(req ports.ZoneRequest, ctx context.Context) (ports.ZoneResponse, error) {
	if req.Zone == "invalid"  {
		return ports.ZoneResponse{}, ports.ErrInvalidInput
	}
	return ports.ZoneResponse{
		Zone:	req.Zone,
	}, nil
}

func (s *ConsumerService) Login(req ports.ConsumerLoginRequest, ctx context.Context) (ports.LoginResponse, error) {
	if req.Username == "alice" && req.Password == "pw" {
		return ports.LoginResponse{Secret: "login-token"}, nil
	}
	return ports.LoginResponse{}, ports.ErrUnauthorized
}

func (s *ConsumerService) Register(req ports.ConsumerRegistrationRequest, ctx context.Context) (ports.RegisterResponse, error) {
	if req.Username == "" || req.Password == "" {
		return ports.RegisterResponse{}, ports.ErrInvalidInput
	}
	return ports.RegisterResponse{Secret: "registered-token"}, nil
}

var _ ports.Api = (*ConsumerService)(nil)
