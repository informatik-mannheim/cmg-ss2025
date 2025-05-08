package core_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/informatik-mannheim/cmg-ss2025/services/job/core"
	"github.com/informatik-mannheim/cmg-ss2025/services/job/ports"
)

// MockJobStorage is a mock implementation of the JobStorage interface
type MockJobStorage struct {
	jobs map[string]ports.Job
}

func NewMockJobStorage() *MockJobStorage {
	return &MockJobStorage{
		jobs: make(map[string]ports.Job),
	}
}

func (m *MockJobStorage) GetJobs(ctx context.Context, status []ports.JobStatus) ([]ports.Job, error) {
	var results []ports.Job
	for _, job := range m.jobs {
		if len(status) == 0 || containsStatus(status, job.Status) {
			results = append(results, job)
		}
	}
	return results, nil
}

func (m *MockJobStorage) CreateJob(ctx context.Context, job ports.Job) error {
	m.jobs[job.Id] = job
	return nil
}

func (m *MockJobStorage) GetJob(ctx context.Context, id string) (ports.Job, error) {
	job, ok := m.jobs[id]
	if !ok {
		return ports.Job{}, errors.New("job not found")
	}
	return job, nil
}

func (m *MockJobStorage) UpdateJob(ctx context.Context, id string, updateFn func(*ports.Job) error) (ports.Job, error) {
	job, ok := m.jobs[id]
	if !ok {
		return ports.Job{}, errors.New("job not found")
	}

	if err := updateFn(&job); err != nil {
		return ports.Job{}, err
	}
	m.jobs[id] = job
	return job, nil
}

func containsStatus(statusList []ports.JobStatus, status ports.JobStatus) bool {
	for _, s := range statusList {
		if s == status {
			return true
		}
	}
	return false
}

func setup() *core.JobService {
	storage := NewMockJobStorage()
	return core.NewJobService(storage)
}

func TestJobService_CreateJob(t *testing.T) {
	service := setup()
	ctx := context.Background()

	tests := []struct {
		name    string
		args    ports.JobCreate
		wantErr bool
	}{
		{
			name: "Create a valid job",
			args: ports.JobCreate{
				JobName:      "Test Job",
				CreationZone: "DE",
				Image:        ports.ContainerImage{Name: "golang", Version: "1.15"},
				Parameters: map[string]string{
					"volumes": "/host/path:/container/path",
					"ports":   "80:8080",
					"env":     "NODE_ENV=development",
				},
			},
			wantErr: false,
		},
		{
			name: "Missing job name",
			args: ports.JobCreate{
				CreationZone: "DE",
				Image:        ports.ContainerImage{Name: "golang", Version: "1.15"},
				Parameters: map[string]string{
					"volumes": "/host/path:/container/path",
					"ports":   "80:8080",
					"env":     "NODE_ENV=development",
				},
			},
			wantErr: true, // Assuming your implementation returns an error
		},
		{
			name: "Invalid image version format",
			args: ports.JobCreate{
				JobName:      "Test Job",
				CreationZone: "DE",
				Image:        ports.ContainerImage{Name: "golang", Version: "!!15"},
				Parameters: map[string]string{
					"volumes": "/host/path:/container/path",
					"ports":   "80:8080",
					"env":     "NODE_ENV=development",
				},
			},
			wantErr: true, // Assuming your implementation returns an error
		},
		{
			name: "Job with special characters",
			args: ports.JobCreate{
				JobName:      "Test @Job!",
				CreationZone: "DE",
				Image:        ports.ContainerImage{Name: "Python", Version: "3.8"},
				Parameters: map[string]string{
					"volumes": "/host/path:/container/path",
					"ports":   "80:8080",
					"env":     "NODE_ENV=development",
				},
			},
			wantErr: false,
		},
		{
			name: "Large parameter list",
			args: ports.JobCreate{
				JobName:      "Large Parameter Job",
				CreationZone: "DE",
				Image:        ports.ContainerImage{Name: "node", Version: "14"},
				Parameters:   generateLargeParameters(1000),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			job, err := service.CreateJob(ctx, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateJob() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if job.Id == "" {
					t.Error("Expected job.Id to be non-empty")
				}
				if job.JobName != tt.args.JobName {
					t.Errorf("Expected job.JobName = %v, got %v", tt.args.JobName, job.JobName)
				}
				if job.Image.Name != tt.args.Image.Name || job.Image.Version != tt.args.Image.Version {
					t.Errorf("Expected job.Image = %v, got %v", tt.args.Image, job.Image)
				}
			}
		})
	}
}

func TestJobService_GetJob(t *testing.T) {
	service := setup()
	ctx := context.Background()

	jobCreate := ports.JobCreate{
		JobName:      "Retrieve Test Job",
		CreationZone: "US",
		Image:        ports.ContainerImage{Name: "node", Version: "14"},
		Parameters: map[string]string{
			"volumes": "/host/path:/container/path",
			"ports":   "80:8080",
			"env":     "NODE_ENV=development",
		},
	}
	createdJob, _ := service.CreateJob(ctx, jobCreate)

	tests := []struct {
		name    string
		id      string
		want    ports.Job
		wantErr bool
	}{
		{
			name:    "Get existing job",
			id:      createdJob.Id,
			want:    createdJob,
			wantErr: false,
		},
		{
			name:    "Get non-existing job",
			id:      uuid.NewString(),
			want:    ports.Job{},
			wantErr: true,
		},
		{
			name:    "Invalid ID format",
			id:      "invalid-uuid",
			want:    ports.Job{},
			wantErr: true,
		},
		{
			name:    "Empty ID",
			id:      "",
			want:    ports.Job{},
			wantErr: true,
		},
		{
			name:    "ID with only spaces",
			id:      "    ",
			want:    ports.Job{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.GetJob(ctx, tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetJob() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				// Only check field values if no error is expected
				if got.Id != tt.want.Id ||
					got.JobName != tt.want.JobName ||
					got.Image.Name != tt.want.Image.Name ||
					got.Image.Version != tt.want.Image.Version {
					t.Errorf("GetJob() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestJobService_GetJobs(t *testing.T) {
	service := setup()
	ctx := context.Background()

	// Create jobs with different statuses
	statuses := []ports.JobStatus{ports.StatusQueued, ports.StatusRunning, ports.StatusCompleted}
	for _, status := range statuses {
		job := ports.JobCreate{
			JobName:      "Test Job " + string(status),
			CreationZone: "DE",
			Image:        ports.ContainerImage{Name: "golang", Version: "1.15"},
			Parameters: map[string]string{
				"volumes": "/host/path:/container/path",
				"ports":   "80:8080",
				"env":     "NODE_ENV=development",
			},
		}
		createdJob, _ := service.CreateJob(ctx, job)

		updateData := ports.SchedulerUpdateData{
			WorkerID:        uuid.NewString(),
			ComputeZone:     "DE",
			CarbonIntensity: 50,
			CarbonSaving:    10,
			Status:          status,
		}

		if _, err := service.UpdateJobScheduler(ctx, createdJob.Id, updateData); err != nil {
			t.Fatalf("Failed to update job: %v", err)
		}
	}

	tests := []struct {
		name    string
		status  []ports.JobStatus
		wantLen int
		wantErr bool
		setup   func() *core.JobService // Function to setup the required scenario
	}{
		{
			name:    "Get all queued jobs",
			status:  []ports.JobStatus{ports.StatusQueued},
			wantLen: 1,
			wantErr: false,
			setup:   func() *core.JobService { return service },
		},
		{
			name:    "Get all jobs without filter",
			status:  nil,
			wantLen: 3,
			wantErr: false,
			setup:   func() *core.JobService { return service },
		},
		{
			name:    "Get jobs with non-existing status",
			status:  []ports.JobStatus{"non-existing-status"},
			wantLen: 0,
			wantErr: false,
			setup:   func() *core.JobService { return service },
		},
		{
			name:    "Get jobs with multiple statuses",
			status:  []ports.JobStatus{ports.StatusQueued, ports.StatusRunning},
			wantLen: 2,
			wantErr: false,
			setup:   func() *core.JobService { return service },
		},
		{
			name:    "Get jobs from an empty storage",
			status:  []ports.JobStatus{ports.StatusCompleted},
			wantLen: 0,
			wantErr: false,
			setup: func() *core.JobService {
				// Create a new service with an empty storage for this case
				return core.NewJobService(NewMockJobStorage())
			},
		},
		{
			name:    "Get jobs with invalid status entry",
			status:  []ports.JobStatus{""},
			wantLen: 0,
			wantErr: false,
			setup:   func() *core.JobService { return service },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := tt.setup()
			got, err := service.GetJobs(ctx, tt.status)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetJobs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.wantLen {
				t.Errorf("GetJobs() got %v jobs, want %v", len(got), tt.wantLen)
			}
		})
	}
}

func TestJobService_GetJobOutcome(t *testing.T) {
	service := setup()
	ctx := context.Background()

	jobCreate := ports.JobCreate{
		JobName:      "Outcome Test Job",
		CreationZone: "EU",
		Image:        ports.ContainerImage{Name: "python", Version: "3.9"},
		Parameters: map[string]string{
			"volumes": "/host/path:/container/path",
			"ports":   "80:8080",
			"env":     "NODE_ENV=development",
		},
	}
	createdJob, _ := service.CreateJob(ctx, jobCreate)

	updateData := ports.WorkerDaemonUpdateData{
		Status:       ports.StatusCompleted,
		Result:       "Analysis complete. Results stored in /data/analysis/output.txt.",
		ErrorMessage: "",
	}
	service.UpdateJobWorkerDaemon(ctx, createdJob.Id, updateData) // Ensure job outcome data is updated

	// Expected outcome setup
	expectedOutcome := ports.JobOutcome{
		JobName:         "Outcome Test Job",
		Status:          "completed",
		Result:          "Analysis complete. Results stored in /data/analysis/output.txt.",
		ErrorMessage:    "",
		ComputeZone:     "", // Assuming not updated, left empty
		CarbonIntensity: 0,  // Assuming default value
		CarbonSavings:   0,  // Assuming default value
	}

	tests := []struct {
		name    string
		id      string
		want    ports.JobOutcome
		wantErr bool
	}{
		{
			name:    "Get outcome of existing job",
			id:      createdJob.Id,
			want:    expectedOutcome,
			wantErr: false,
		},
		{
			name:    "Get outcome of non-existing job",
			id:      uuid.NewString(),
			want:    ports.JobOutcome{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.GetJobOutcome(ctx, tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetJobOutcome() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got.JobName != tt.want.JobName ||
				got.Status != tt.want.Status ||
				got.Result != tt.want.Result ||
				got.ErrorMessage != tt.want.ErrorMessage ||
				got.ComputeZone != tt.want.ComputeZone ||
				got.CarbonIntensity != tt.want.CarbonIntensity ||
				got.CarbonSavings != tt.want.CarbonSavings {
				t.Errorf("GetJobOutcome() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJobService_UpdateJobScheduler(t *testing.T) {
	service := setup()
	ctx := context.Background()

	jobCreate := ports.JobCreate{
		JobName:      "Update Scheduler Test",
		CreationZone: "FR",
		Image:        ports.ContainerImage{Name: "python", Version: "3.8"},
		Parameters:   map[string]string{"GPU": "NVIDIA"},
	}
	createdJob, _ := service.CreateJob(ctx, jobCreate)
	updateData := ports.SchedulerUpdateData{
		WorkerID:        uuid.NewString(),
		ComputeZone:     "FR",
		CarbonIntensity: 75,
		CarbonSaving:    30,
		Status:          ports.StatusScheduled,
	}

	tests := []struct {
		name    string
		id      string
		data    ports.SchedulerUpdateData
		wantErr bool
	}{
		{
			name:    "Update existing job",
			id:      createdJob.Id,
			data:    updateData,
			wantErr: false,
		},
		{
			name:    "Update non-existing job",
			id:      uuid.NewString(),
			data:    updateData,
			wantErr: true,
		},
		{
			name: "Missing worker ID",
			id:   createdJob.Id,
			data: ports.SchedulerUpdateData{
				WorkerID:        "", // Missing worker ID
				ComputeZone:     "FR",
				CarbonIntensity: 75,
				CarbonSaving:    30,
				Status:          ports.StatusScheduled,
			},
			wantErr: true,
		},
		{
			name: "Negative carbon intensity",
			id:   createdJob.Id,
			data: ports.SchedulerUpdateData{
				WorkerID:        uuid.NewString(),
				ComputeZone:     "FR",
				CarbonIntensity: -1, // Negative value
				CarbonSaving:    30,
				Status:          ports.StatusScheduled,
			},
			wantErr: true,
		},
		{
			name:    "Empty ID",
			id:      "",
			data:    updateData,
			wantErr: true,
		},
		{
			name:    "Invalid ID format",
			id:      "not-a-uuid",
			data:    updateData,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.UpdateJobScheduler(ctx, tt.id, tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateJobScheduler() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestJobService_UpdateJobWorkerDaemon(t *testing.T) {
	service := setup()
	ctx := context.Background()

	jobCreate := ports.JobCreate{
		JobName:      "Update Worker Daemon Test",
		CreationZone: "GB",
		Image:        ports.ContainerImage{Name: "java", Version: "8"},
		Parameters:   map[string]string{"Memory": "16GB"},
	}
	createdJob, _ := service.CreateJob(ctx, jobCreate)

	validUpdateData := ports.WorkerDaemonUpdateData{
		Status:       ports.StatusCompleted,
		Result:       "Job completed successfully.",
		ErrorMessage: "",
	}

	tests := []struct {
		name    string
		id      string
		data    ports.WorkerDaemonUpdateData
		wantErr bool
	}{
		{
			name:    "Update existing job",
			id:      createdJob.Id,
			data:    validUpdateData,
			wantErr: false,
		},
		{
			name:    "Update non-existing job",
			id:      uuid.NewString(),
			data:    validUpdateData,
			wantErr: true,
		},
		{
			name:    "Invalid ID format",
			id:      "not-a-uuid",
			data:    validUpdateData,
			wantErr: true,
		},
		{
			name:    "Empty ID",
			id:      "",
			data:    validUpdateData,
			wantErr: true,
		},
		{
			name:    "Empty Result and ErrorMessage",
			id:      createdJob.Id,
			data:    ports.WorkerDaemonUpdateData{Status: ports.StatusCompleted, Result: "", ErrorMessage: ""},
			wantErr: false,
		},
		{
			name: "Failed status with ErrorMessage",
			id:   createdJob.Id,
			data: ports.WorkerDaemonUpdateData{
				Status:       ports.StatusFailed,
				Result:       "",
				ErrorMessage: "Execution error occurred.",
			},
			wantErr: false,
		},
		{
			name: "Scheduled status without Result",
			id:   createdJob.Id,
			data: ports.WorkerDaemonUpdateData{
				Status:       ports.StatusScheduled,
				Result:       "",
				ErrorMessage: "",
			},
			wantErr: false, // Assuming this is acceptable; adjust according to your logic
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.UpdateJobWorkerDaemon(ctx, tt.id, tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateJobWorkerDaemon() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func generateLargeParameters(n int) map[string]string {
	params := make(map[string]string)
	for i := 0; i < n; i++ {
		params[fmt.Sprintf("param%d", i)] = fmt.Sprintf("value%d", i)
	}
	return params
}
