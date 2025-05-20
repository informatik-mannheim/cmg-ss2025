package core_test

import (
	"context"
	"errors"
	"testing"

	"github.com/informatik-mannheim/cmg-ss2025/services/consumer-gateway/core"
	"github.com/informatik-mannheim/cmg-ss2025/services/consumer-gateway/ports"
)

type mockJobClient struct {
	createJobCalled  bool
	getOutcomeCalled bool
	failCreate       bool
	failOutcome      bool
}

func (m *mockJobClient) CreateJob(ctx context.Context, req ports.CreateJobRequest) (ports.CreateJobResponse, error) {
	m.createJobCalled = true
	if m.failCreate {
		return ports.CreateJobResponse{}, ports.ErrInvalidInput
	}
	return ports.CreateJobResponse{
		ImageID: req.ImageID,
		JobName: req.JobName,
	}, nil
}

func (m *mockJobClient) GetJobOutcome(ctx context.Context, jobID string) (ports.JobOutcomeResponse, error) {
	m.getOutcomeCalled = true
	if m.failOutcome {
		return ports.JobOutcomeResponse{}, ports.ErrNotFound
	}
	return ports.JobOutcomeResponse{
		JobName: "job-1", Status: ports.JobStatus("done"), Result: "http://example.com/image.png",
	}, nil
}

type mockZoneClient struct {
	fail bool
}

func (m *mockZoneClient) GetZone(ctx context.Context, req ports.ZoneRequest) (ports.ZoneResponse, error) {
	if m.fail {
		return ports.ZoneResponse{}, ports.ErrBadRequest
	}
	return ports.ZoneResponse{Zone: req.Zone}, nil
}

type mockLoginClient struct {
	fail bool
}

func (m *mockLoginClient) Login(ctx context.Context, req ports.ConsumerLoginRequest) (ports.LoginResponse, error) {
	if m.fail {
		return ports.LoginResponse{}, ports.ErrUnauthorized
	}
	return ports.LoginResponse{Secret: "token-123"}, nil
}

func TestConsumerGatewayService_CreateJob(t *testing.T) {
	jobMock := &mockJobClient{}
	service := core.NewConsumerService(jobMock, &mockZoneClient{}, &mockLoginClient{})

	resp, err := service.CreateJob(context.Background(), ports.CreateJobRequest{
		ImageID:      "img1",
		JobName:      "job-1",
		CreationZone: "GER",
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.ImageID != "img1" {
		t.Errorf("unexpected ImageID: %v", resp.ImageID)
	}
}

func TestConsumerGatewayService_GetJobOutcome(t *testing.T) {
	jobMock := &mockJobClient{}
	service := core.NewConsumerService(jobMock, &mockZoneClient{}, &mockLoginClient{})

	resp, err := service.GetJobOutcome(context.Background(), "job-123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Status != "done" {
		t.Errorf("unexpected status: %v", resp.Status)
	}
}

func TestConsumerGatewayService_GetZone(t *testing.T) {
	service := core.NewConsumerService(&mockJobClient{}, &mockZoneClient{}, &mockLoginClient{})

	resp, err := service.GetZone(context.Background(), ports.ZoneRequest{Zone: "GER"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Zone != "GER" {
		t.Errorf("unexpected zone: %v", resp.Zone)
	}
}

func TestConsumerGatewayService_Login(t *testing.T) {
	service := core.NewConsumerService(&mockJobClient{}, &mockZoneClient{}, &mockLoginClient{})

	resp, err := service.Login(context.Background(), ports.ConsumerLoginRequest{
		Username: "alice",
		Password: "pw",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Secret != "token-123" {
		t.Errorf("unexpected secret: %v", resp.Secret)
	}
}

func TestConsumerGatewayService_Login_Unauthorized(t *testing.T) {
	service := core.NewConsumerService(&mockJobClient{}, &mockZoneClient{}, &mockLoginClient{fail: true})

	_, err := service.Login(context.Background(), ports.ConsumerLoginRequest{
		Username: "wrong",
		Password: "bad",
	})
	if !errors.Is(err, ports.ErrUnauthorized) {
		t.Errorf("expected ErrUnauthorized, got %v", err)
	}
}
