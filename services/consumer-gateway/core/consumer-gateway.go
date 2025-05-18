package core

import (
	"context"

	"github.com/informatik-mannheim/cmg-ss2025/services/consumer-gateway/ports"
)

type ConsumerGatewayService struct {
	job   ports.JobClient
	zone  ports.CarbonIntensityClient
	login ports.LoginClient
}

func NewConsumerService(jobClient ports.JobClient, zoneClient ports.CarbonIntensityClient, loginClient ports.LoginClient) *ConsumerGatewayService {
	return &ConsumerGatewayService{
		job:   jobClient,
		zone:  zoneClient,
		login: loginClient,
	}
}

func (s *ConsumerGatewayService) CreateJob(ctx context.Context, req ports.CreateJobRequest) (ports.CreateJobResponse, error) {
	return s.job.CreateJob(ctx, req)
}

func (s *ConsumerGatewayService) GetJobOutcome(ctx context.Context, jobID string) (ports.JobOutcomeResponse, error) {
	return s.job.GetJobOutcome(ctx, jobID)
}

func (s *ConsumerGatewayService) GetZone(req ports.ZoneRequest, ctx context.Context) (ports.ZoneResponse, error) {
	return s.zone.GetZone(req, ctx)
}

func (s *ConsumerGatewayService) Login(req ports.ConsumerLoginRequest, ctx context.Context) (ports.LoginResponse, error) {
	return s.login.Login(req, ctx)
}

var _ ports.JobClient = (*ConsumerGatewayService)(nil)
