package core_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"


	"github.com/informatik-mannheim/cmg-ss2025/services/consumer-gateway/ports"
	"github.com/informatik-mannheim/cmg-ss2025/services/consumer-gateway/adapters/handler-http"
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
		Secret : "secret-123"}, ports.ErrUnauthorized
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
		ImageID:     "job-123",
		JobStatus: "queued",
	}, nil
}


func (f *FakeService) GetJobResult(_ string, ctx context.Context) (ports.JobResultResponse, error) {
	user := ctx.Value(userContextKey).(string)
	if user == "alice" {
		return ports.JobResultResponse{
			ImageID: "job-123",
			JobStatus: "completed",
		}, nil
	}
	return ports.JobResultResponse{}, ports.ErrNotFound
}


/*
Will be implemented, once the /zones endpoint is on main.
func (f *FakeService) GetZones(ctx context.Context) (ports.ZonesResponse, error) {
	return ports.ZonesResponse{}, nil // wenn nötig
}*/ 



func TestConsumerService(t *testing.T) {
	
	 tests := [] struct {
		user 			string
		name   			string
		url 			string
		method			string
		code   			int
		inputBody		string
		responseBody 	string


	}{ 		// Login Tests
		{ 
			name: "Successful Login",
			url: "/auth/login",
			method: "POST",
			inputBody: `{"username" : "Alice Bob", "password": "SuperSecure123"}`,
			responseBody: `{"secret": "super-secret-123"}`,
			code: http.StatusOK,

		},
		{
			name: "Unsuccessful Login",
			url: "/auth/login",
			method: "POST",
			inputBody: `{"username" : "Bad User", "password": "wrong"}`,
			responseBody: `{"error": "Unauthorized}`,
			code: http.StatusBadRequest,
		},
		{
			name: "Invalid JSON Login",
			url: "/auth/login",
			method: "POST",
			inputBody: `{"username" : "", "password": ""}`,
			responseBody: `{"error": "Invalid Request}`,
			code: http.StatusBadRequest,			
		},
		

			// Registration Tests
		{
			name: "Successful registration",
			url: "/auth/register",
			method: "POST",
			code: http.StatusOK,
			inputBody: `{"username" : "Alice Bob", "password": "SuperSecure123"}`,
			responseBody: `{"secret": "super-secret-123"}`,
		},
		{
			name: "Unsuccessful registration",
			url: "/auth/register",
			method: "POST",
			code: http.StatusUnauthorized,
			inputBody: `{"username" : "Bad User", "password": "wrong"}`,
			responseBody: `{"error": "Unauthorized"}`,
		},
		{
			name: "Invalid JSON Registration",
			url: "/auth/register",
			method: "POST",
			code: http.StatusUnauthorized,
			inputBody: `{"username" : "", "password": ""}`,
			responseBody: `{"error": "Invalid Request"}`,
		},
			// Create Job test
		{
			name: "Create Job successfully",
			url: "/jobs",
			method: "POST",
			code: http.StatusOK,
			inputBody: `{"image_id": "123", "Zone" : "FR", "params":  "abc"}`,
			responseBody: `{"image_id:" "job-123", "job_status": "queued"}`,

		},
		{
			name: "Create Job – invalid ID",
			url: "/jobs",
			method: "POST",
			code: http.StatusUnauthorized,
			inputBody: `{"image_id": "invalid", "zone" : "FR", "params":  "abc"}`,
			responseBody: `{"error": "Invalid Request"}`,

		},
		{
			name: "Create Job – invalid zone",
			url: "/jobs",
			method: "POST",
			code: http.StatusUnauthorized,
			inputBody: `{"image_id": "123", "zone" : "invalid", "params":  "abc"}`,
			responseBody: `{"error": "Invalid Request"}`,

		},

		{
			name: "Create Job – wrong params",
			url: "/jobs",
			method: "POST",
			code: http.StatusUnauthorized,
			inputBody: `{"image_id": "123", "zone" : "FR", "params":  "invalid"}`,
			responseBody: `{"error": "Invalid Request"}`,

		},
			// Get Job result tests
		{
			name: "Get a job result",
			url: "/jobs/job-123/result",
			method: "GET",
			code: http.StatusOK,
			user: "alice",
			inputBody: ``,
			responseBody: `{"image_id": "job-123", "job_status": "completed"}`,
		},
		{
			name: "Job not found",
			url: "/jobs/{id}/results",
			method: "GET",
			code: http.StatusNotFound,
			inputBody: `{"image_id": "wrong"}`,
			responseBody: `{"error": "Invalid Request"}`,
		},
	}


		
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &FakeService{}
			router := handler_http.NewHandler(service)

			body := strings.NewReader(tt.inputBody)
			req := httptest.NewRequest(tt.method, tt.url, body)
			req.Header.Set("Content-Type", "application/json")

			// This simulates the JWT Token in the header
			req.Header.Set("Authorization", "Bearer alice-token")
			ctx := context.WithValue(req.Context(), userContextKey, "alice")
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			if rr.Code != tt.code {
				t.Errorf("expected code %d, got %d", tt.code, rr.Code)
				}
	})
}
}