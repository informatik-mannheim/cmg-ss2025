package core

import (
	"context"
	"github.com/informatik-mannheim/cmg-ss2025/services/consumer-gateway/ports"
)

type ConsumerGatewayService struct {
	job   ports.JobClient
	zone  ports.ZoneClient
	login ports.LoginClient
}

var _ ports.Api = &ConsumerGatewayService{}

func NewConsumerService(jobClient ports.JobClient, zoneClient ports.ZoneClient, loginClient ports.LoginClient) *ConsumerGatewayService {
	return &ConsumerGatewayService{
		job:   jobClient,
		zone:  zoneClient,
		login: loginClient,
	}
}

func (s *ConsumerGatewayService) CreateJob(ctx context.Context, req ports.CreateJobRequest) (ports.CreateJobResponse, error) {
	resp, err := s.job.CreateJob(ctx, req)
	if err != nil {
		return ports.CreateJobResponse{}, err
	}
	return resp, nil
}

func (s *ConsumerGatewayService) GetJobOutcome(ctx context.Context, jobID string) (ports.JobOutcomeResponse, error) {
	resp, err := s.job.GetJobOutcome(ctx, jobID)
	if err != nil {
		return ports.JobOutcomeResponse{}, err
	}
	return resp, nil
}

func (s *ConsumerGatewayService) GetZone(ctx context.Context, req ports.ZoneRequest) (ports.ZoneResponse, error) {
	resp, err := s.zone.GetZone(ctx, req)
	if err != nil {
		return ports.ZoneResponse{}, err
	}
	return resp, nil
}

func (s *ConsumerGatewayService) Login(ctx context.Context, req ports.ConsumerLoginRequest) (ports.LoginResponse, error) {
	resp, err := s.login.Login(ctx, ports.LoginClientRequest{Secret: req.Secret})
	if err != nil {
		return ports.LoginResponse{}, err
	}
	return ports.LoginResponse{Secret: resp.Token}, nil
}

var _ ports.JobClient = (*ConsumerGatewayService)(nil)
