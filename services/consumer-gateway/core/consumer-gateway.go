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

func (s *ConsumerService) CreateJob(req ports.CreateJobRequest, ctx context.Context)(ports.CreateJobResponse, error){
	return ports.CreateJobResponse{}, nil
}

func (s *ConsumerService) GetJobResult(jobID string, ctx context.Context)(ports.JobResultResponse, error){
	return ports.JobResultResponse{}, nil
}

func (s *ConsumerService) GetZones(req ports.GetZones, ctx context.Context)(ports.ZonesResponse, error){
	return ports.ZonesResponse{}, nil
}

func (s *ConsumerService) Login(req ports.ConsumerLoginRequest, ctx context.Context)(ports.LoginResponse, error) {
	return ports.LoginResponse{}, nil
}

func (s *ConsumerService) Register(req ports.ConsumerRegistrationRequest, ctx context.Context)(ports.RegisterResponse, error){
	return ports.RegisterResponse{}, nil
}


 var _ ports.Api = (*ConsumerService)(nil)
