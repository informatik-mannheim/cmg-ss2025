package handler_http

import (
	"context"

	"github.com/informatik-mannheim/cmg-ss2025/services/consumer-gateway/ports"
)

type contextKey string

const userContextKey contextKey = "user"

type FakeService struct{}

func (f *FakeService) Login(ctx context.Context, req ports.ConsumerLoginRequest) (ports.LoginResponse, error) {
	if req.Username == "Alice Bob" {
		return ports.LoginResponse{Secret: "abc-123"}, nil
	}
	return ports.LoginResponse{}, ports.ErrUnauthorized
}

func (f *FakeService) CreateJob(ctx context.Context, req ports.CreateJobRequest) (ports.CreateJobResponse, error) {
	if req.CreationZone == "" || req.Parameters == nil {
		return ports.CreateJobResponse{}, ports.ErrInvalidInput
	}
	return ports.CreateJobResponse{
		JobName:      "job-123",
		CreationZone: req.CreationZone,
		Parameters:   req.Parameters,
		Status:       "queued",
	}, nil
}

func (f *FakeService) GetJobOutcome(ctx context.Context, jobID string) (ports.JobOutcomeResponse, error) {
	user := ctx.Value(userContextKey).(string)
	if user == "alice" && jobID == "job-123" {
		return ports.JobOutcomeResponse{
			JobName:       "job-123",
			Status:        "completed",
			Result:        "Result",
			CarbonSavings: 123,
		}, nil
	}
	return ports.JobOutcomeResponse{}, ports.ErrNotFound
}

func (f *FakeService) GetZone(ctx context.Context, req ports.ZoneRequest) (ports.ZoneResponse, error) {
	if req.Zone == "" || req.Zone == "invalid" {
		return ports.ZoneResponse{}, ports.ErrInvalidInput
	}
	return ports.ZoneResponse{Zone: req.Zone}, nil
}
