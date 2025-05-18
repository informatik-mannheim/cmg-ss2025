package handler_http

import (
	"context"

	"github.com/informatik-mannheim/cmg-ss2025/services/consumer-gateway/ports"
)

// This is to avoid using a string in the context
type contextKey string

const userContextKey contextKey = "user"

// The following functions return mock data
type FakeService struct{}

func (f *FakeService) Login(req ports.ConsumerLoginRequest, ctx context.Context) (ports.LoginResponse, error) {
	if req.Username == "Alice Bob" {
		return ports.LoginResponse{Secret: "abc-123"}, nil
	}
	return ports.LoginResponse{}, ports.ErrUnauthorized
}

func (f *FakeService) Register(req ports.ConsumerRegistrationRequest, ctx context.Context) (ports.RegisterResponse, error) {
	if req.Username == "Alice Bob" {
		return ports.RegisterResponse{Secret: "abc-123"}, nil
	}
	if req.Username == "invalid" || req.Username == "" {
		return ports.RegisterResponse{}, ports.ErrInvalidInput
	}
	return ports.RegisterResponse{
		Secret: "secret-123"}, ports.ErrUnauthorized
}

func (f *FakeService) CreateJob(req ports.CreateJobRequest, ctx context.Context) (ports.CreateJobResponse, error) {
	if req.ImageID == "" || req.ImageID == "invalid" {
		return ports.CreateJobResponse{}, ports.ErrInvalidInput
	}
	if req.Zone == "" || req.Zone == "invalid" {
		return ports.CreateJobResponse{}, ports.ErrInvalidInput
	}
	// Is allowed to be empty, but not invalid
	if req.Param == "invalid" {
		return ports.CreateJobResponse{}, ports.ErrInvalidInput
	}
	return ports.CreateJobResponse{
		ImageID:   "job-123",
		JobStatus: "queued",
	}, nil
}

func (f *FakeService) GetJobResult(_ string, ctx context.Context) (ports.JobOutcomeResponse, error) {
	user := ctx.Value(userContextKey).(string)
	if user == "alice" {
		return ports.JobOutcomeResponse{
			ImageID:   "job-123",
			JobStatus: "completed",
		}, nil
	}
	return ports.JobOutcomeResponse{}, ports.ErrNotFound
}

/*
Will be implemented, once the /zones endpoint is on main.
func (f *FakeService) GetZones(ctx context.Context) (ports.ZonesResponse, error) {
	return ports.ZonesResponse{}, nil // wenn nötig
}*/
