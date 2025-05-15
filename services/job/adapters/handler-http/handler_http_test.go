package handler_http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	handler_http "github.com/informatik-mannheim/cmg-ss2025/services/job/adapters/handler-http"
	"github.com/informatik-mannheim/cmg-ss2025/services/job/ports"
)

// MockJobService implements the JobService interface for testing purposes.
type MockJobService struct{}

func (m *MockJobService) GetJobs(_ context.Context, status []ports.JobStatus) ([]ports.Job, error) {
	if len(status) == 0 {
		return nil, nil
	}
	return []ports.Job{{Id: "123", JobName: "mockJob"}}, nil
}

func (m *MockJobService) CreateJob(_ context.Context, jobCreate ports.JobCreate) (ports.Job, error) {
	if jobCreate.JobName == "" {
		return ports.Job{}, ports.ErrNotExistingJobName
	}
	return ports.Job{Id: "123", JobName: jobCreate.JobName}, nil
}

func (m *MockJobService) GetJob(_ context.Context, id string) (ports.Job, error) {
	if id == "123" {
		return ports.Job{Id: "123"}, nil
	}
	return ports.Job{}, ports.ErrJobNotFound
}

func (m *MockJobService) GetJobOutcome(_ context.Context, id string) (ports.JobOutcome, error) {
	if id == "123" {
		return ports.JobOutcome{JobName: "mockJobOutcome"}, nil
	}
	return ports.JobOutcome{}, ports.ErrJobNotFound
}

func (m *MockJobService) UpdateJobScheduler(_ context.Context, id string, data ports.SchedulerUpdateData) (ports.Job, error) {
	if id == "123" {
		return ports.Job{Id: "123"}, nil
	}
	return ports.Job{}, ports.ErrJobNotFound
}

func (m *MockJobService) UpdateJobWorkerDaemon(_ context.Context, id string, data ports.WorkerDaemonUpdateData) (ports.Job, error) {
	if id == "123" {
		return ports.Job{Id: "123"}, nil
	}
	return ports.Job{}, ports.ErrJobNotFound
}

func TestHandler_GetJobs(t *testing.T) {
	mockService := &MockJobService{}
	handler := handler_http.NewHandler(mockService)

	tests := []struct {
		name           string
		query          string
		expectedStatus int
	}{
		{"No Status Filter", "", http.StatusNoContent},
		{"Valid Status", "?status=queued", http.StatusOK},
		{"Invalid Status", "?status=invalid", http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/jobs"+tt.query, nil)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %v; got %v", tt.expectedStatus, rr.Code)
			}
		})
	}
}

func TestHandler_CreateJob(t *testing.T) {
	mockService := &MockJobService{}
	handler := handler_http.NewHandler(mockService)

	tests := []struct {
		name           string
		payload        string
		expectedStatus int
	}{
		{"Valid Job", `{"JobName":"New Job", "CreationZone":"DE", "Image":{"Name":"test-image", "Version":"1.0"}}`, http.StatusCreated},
		{"Empty Job Name", `{"JobName":"", "CreationZone":"DE", "Image":{"Name":"test-image", "Version":"1.0"}}`, http.StatusBadRequest},
		{"Invalid JSON", `invalid-json`, http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/jobs", strings.NewReader(tt.payload))
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %v; got %v", tt.expectedStatus, rr.Code)
			}
		})
	}
}

func TestHandler_GetJob(t *testing.T) {
	mockService := &MockJobService{}
	handler := handler_http.NewHandler(mockService)

	tests := []struct {
		name           string
		id             string
		expectedStatus int
	}{
		{"Existing Job", "123", http.StatusOK},
		{"Non-Existing Job", "456", http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/jobs/"+tt.id, nil)
			rr := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/jobs/{id}", handler.GetJob)
			router.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %v; got %v", tt.expectedStatus, rr.Code)
			}
		})
	}
}

func TestHandler_GetJobOutcome(t *testing.T) {
	mockService := &MockJobService{}
	handler := handler_http.NewHandler(mockService)

	tests := []struct {
		name           string
		id             string
		expectedStatus int
	}{
		{"Existing Job Outcome", "123", http.StatusOK},
		{"Non-Existing Job Outcome", "456", http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/jobs/"+tt.id+"/outcome", nil)
			rr := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/jobs/{id}/outcome", handler.GetJobOutcome)
			router.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %v; got %v", tt.expectedStatus, rr.Code)
			}
		})
	}
}

func TestHandler_UpdateJobScheduler(t *testing.T) {
	mockService := &MockJobService{}
	handler := handler_http.NewHandler(mockService)

	updateData := ports.SchedulerUpdateData{
		WorkerID:        "worker-id",
		ComputeZone:     "DE",
		CarbonIntensity: 75,
		CarbonSaving:    5,
		Status:          ports.StatusScheduled,
	}
	payload, _ := json.Marshal(updateData)

	req, _ := http.NewRequest("PATCH", "/jobs/123/update-scheduler", bytes.NewBuffer(payload))
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/jobs/{id}/update-scheduler", handler.UpdateJobScheduler)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestHandler_UpdateJobWorkerDaemon(t *testing.T) {
	mockService := &MockJobService{}
	handler := handler_http.NewHandler(mockService)

	updateData := ports.WorkerDaemonUpdateData{
		Status:       ports.StatusCompleted,
		Result:       "success",
		ErrorMessage: "",
	}
	payload, _ := json.Marshal(updateData)

	req, _ := http.NewRequest("PATCH", "/jobs/123/update-workerdaemon", bytes.NewBuffer(payload))
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/jobs/{id}/update-workerdaemon", handler.UpdateJobWorkerDaemon)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}
