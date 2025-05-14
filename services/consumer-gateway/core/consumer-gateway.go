package core

import (
	"context"
	"errors"

	"github.com/informatik-mannheim/cmg-ss2025/services/consumer-gateway/ports"
)

type ConsumerService struct {

}

func NewConsumerService() *ConsumerService {
	return &ConsumerService{}
}

func (s *ConsumerService) CreateJob(req ports.CreateJobRequest, ctx context.Context)(ports.CreateJobResponse, error){
	return ports.CreateJobResponse{}, errors.New("error creating job")
}

func (s *ConsumerService) GetJobResult(jobID string, ctx context.Context)(ports.JobResultResponse, error){
	return ports.JobResultResponse{}, errors.New("error getting job result")
}
func (s *ConsumerService) Login(req ports.ConsumerLoginRequest, ctx context.Context)(ports.LoginResponse, error) {
	return ports.LoginResponse{}, errors.New("error during login")
}

func (s *ConsumerService) Register(req ports.ConsumerRegistrationRequest, ctx context.Context)(ports.RegisterResponse, error){
	return ports.RegisterResponse{}, errors.New("error during registration")
}


 var _ ports.Api = (*ConsumerService)(nil)
