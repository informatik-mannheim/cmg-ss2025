package core_test

import (
	"context"
	"testing"

	"github.com/informatik-mannheim/cmg-ss2025/services/consumer-gateway/core"
	"github.com/informatik-mannheim/cmg-ss2025/services/consumer-gateway/ports"
)

func TestCreateJob(t *testing.T) {
	service := core.NewConsumerService()

	tests := []struct {
		name    string
		req     ports.CreateJobRequest
		wantErr error
	}{
		{
			name: "valid job",
			req:  ports.CreateJobRequest{ImageID: "img1", Zone: "EU", Param: "-x"},
			wantErr: nil,
		},
		{
			name: "missing image_id",
			req:  ports.CreateJobRequest{ImageID: "", Zone: "EU", Param: "-x"},
			wantErr: ports.ErrInvalidInput,
		},
		{
			name: "invalid param",
			req:  ports.CreateJobRequest{ImageID: "img1", Zone: "EU", Param: "invalid"},
			wantErr: ports.ErrInvalidInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.CreateJob(tt.req, context.Background())
			if err != tt.wantErr {
				t.Errorf("expected error %v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestRegister(t *testing.T) {
	service := core.NewConsumerService()

	tests := []struct {
		name    string
		req     ports.ConsumerRegistrationRequest
		wantErr error
	}{
		{
			name: "valid",
			req:  ports.ConsumerRegistrationRequest{Username: "alice", Password: "pw"},
			wantErr: nil,
		},
		{
			name: "invalid",
			req:  ports.ConsumerRegistrationRequest{Username: "", Password: ""},
			wantErr: ports.ErrInvalidInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.Register(tt.req, context.Background())
			if err != tt.wantErr {
				t.Errorf("expected %v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	service := core.NewConsumerService()

	tests := []struct {
		name    string
		req     ports.ConsumerLoginRequest
		wantErr error
	}{
		{
			name: "valid login",
			req:  ports.ConsumerLoginRequest{Username: "alice", Password: "pw"},
			wantErr: nil,
		},
		{
			name: "unauthorized",
			req:  ports.ConsumerLoginRequest{Username: "invalid", Password: "wrong"},
			wantErr: ports.ErrUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.Login(tt.req, context.Background())
			if err != tt.wantErr {
				t.Errorf("expected %v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestGetJobResult(t *testing.T) {
	service := core.NewConsumerService()

	t.Run("found", func(t *testing.T) {
		resp, err := service.GetJobResult("job-123", context.WithValue(context.Background(), "user", "alice"))
		if err != nil || resp.ImageID == "" {
			t.Errorf("expected job result, got error: %v", err)
		}
	})

	t.Run("not found", func(t *testing.T) {
		_, err := service.GetJobResult("wrong-id", context.WithValue(context.Background(), "user", "other"))
		if err != ports.ErrNotFound {
			t.Errorf("expected not found, got %v", err)
		}
	})
}
