package core

import (
	"context"
	"errors"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
	"github.com/informatik-mannheim/cmg-ss2025/services/job/ports"
)

// JobService is a struct that implements the JobService interface.
// It provides methods to manage jobs, including creating, retrieving, and updating jobs.
// It uses a storage interface to interact with the underlying data store.
type JobService struct {
	storage ports.JobStorage
}

// NewJobService creates a new instance of JobService with the provided storage.
// An error is returned if an invalid (e.g. nil) JobStorage is transferred.
func NewJobService(storage ports.JobStorage) (*JobService, error) {
	if storage == nil {
		return nil, errors.New("storage cannot be nil")
	}

	return &JobService{
		storage: storage,
	}, nil
}

// GetJobs retrieves jobs based on their status.
// It returns a slice of jobs that match the provided status.
// If no status is provided, it returns all jobs.
func (s *JobService) GetJobs(ctx context.Context, status []ports.JobStatus) ([]ports.Job, error) {
	// If status is nil, return all jobs
	if status == nil {
		return nil, errors.New("atleast one status must be provided")
	}

	// Filter out invalid statuses
	validStatuses := make([]ports.JobStatus, 0, len(status))
	for _, s := range status {
		if isValidStatus(s) {
			validStatuses = append(validStatuses, s)
		}
	}

	// If no valid statuses available, return no jobs
	if len(validStatuses) == 0 {
		return []ports.Job{}, nil
	}

	return s.storage.GetJobs(ctx, validStatuses)
}

// CreateJob creates a new job with the provided job creation data.
// It generates a unique ID for the job and sets the initial status to "queued".
// The job is then stored in the storage.
func (s *JobService) CreateJob(ctx context.Context, jobCreate ports.JobCreate) (ports.Job, error) {
	// Input validation
	if strings.TrimSpace(jobCreate.JobName) == "" {
		return ports.Job{}, errors.New("job name must be provided")
	}
	if strings.TrimSpace(jobCreate.CreationZone) == "" {
		return ports.Job{}, errors.New("creation zone must be provided")
	}
	if strings.TrimSpace(jobCreate.Image.Name) == "" || strings.TrimSpace(jobCreate.Image.Version) == "" {
		return ports.Job{}, errors.New("both image name and version must be provided")
	}

	// Version format validation without regex
	if !isSimpleValidVersion(jobCreate.Image.Version) {
		return ports.Job{}, errors.New("image version format is invalid")
	}

	if len(jobCreate.Parameters) == 0 {
		return ports.Job{}, errors.New("at least one parameter must be provided")
	}
	for key, value := range jobCreate.Parameters {
		if strings.TrimSpace(key) == "" || strings.TrimSpace(value) == "" {
			return ports.Job{}, errors.New("parameters cannot have empty keys or values")
		}
	}

	newJob := ports.Job{
		Id:                   uuid.NewString(),
		UserID:               "some-user-id", // this should be replaced with actual user ID from context
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
		JobName:              jobCreate.JobName,
		Image:                jobCreate.Image,
		AdjustmentParameters: jobCreate.Parameters,
		CreationZone:         jobCreate.CreationZone,
		Status:               ports.StatusQueued,
	}

	err := s.storage.CreateJob(ctx, newJob)
	if err != nil {
		return ports.Job{}, err
	}
	return newJob, nil
}

// GetJob retrieves a specific job by its ID.
// It returns the job if found, or an error if not.
// This method is useful for getting detailed information about a specific job.
func (s *JobService) GetJob(ctx context.Context, id string) (ports.Job, error) {
	// Check for empty or whitespace-only ID
	if len(strings.TrimSpace(id)) == 0 {
		return ports.Job{}, errors.New("job ID must be provided")
	}

	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return ports.Job{}, errors.New("job ID must be a valid UUID")
	}
	return s.storage.GetJob(ctx, id)
}

// GetJobOutcome retrieves the outcome of a job by its ID.
// It returns the job name, status, result, error message, compute zone, carbon intensity, and carbon savings.
// If the job is not found, it returns an error.
// This method is useful for getting detailed information about a job's execution.
func (s *JobService) GetJobOutcome(ctx context.Context, id string) (ports.JobOutcome, error) {
	job, err := s.GetJob(ctx, id)
	if err != nil {
		return ports.JobOutcome{}, err
	}

	return ports.JobOutcome{
		JobName:         job.JobName,
		Status:          string(job.Status),
		Result:          job.Result,
		ErrorMessage:    job.ErrorMessage,
		ComputeZone:     job.ComputeZone,
		CarbonIntensity: job.CarbonIntensity,
		CarbonSavings:   job.CarbonSaving,
	}, nil
}

// UpdateJobScheduler updates the job with the provided ID using the provided scheduler update data.
// It modifies the job's worker ID, compute zone, carbon intensity, carbon savings, and status.
// The updated job is returned.
// functional options are used to modify the job's properties.
func (s *JobService) UpdateJobScheduler(ctx context.Context, id string, data ports.SchedulerUpdateData) (ports.Job, error) {
	// Validate ID
	if len(strings.TrimSpace(id)) == 0 {
		return ports.Job{}, errors.New("job ID must be provided")
	}
	if _, err := uuid.Parse(id); err != nil {
		return ports.Job{}, errors.New("job ID must be a valid UUID")
	}

	// Validate WorkerID
	if len(strings.TrimSpace(data.WorkerID)) == 0 {
		return ports.Job{}, errors.New("worker ID must be provided")
	}

	// Validate CarbonIntensity
	if data.CarbonIntensity < 0 {
		return ports.Job{}, errors.New("carbon intensity must be non-negative")
	}

	updateFn := func(job *ports.Job) error {
		job.WorkerID = data.WorkerID
		job.ComputeZone = data.ComputeZone
		job.CarbonIntensity = data.CarbonIntensity
		job.CarbonSaving = data.CarbonSaving
		job.Status = data.Status
		job.UpdatedAt = time.Now()
		return nil
	}

	return s.storage.UpdateJob(ctx, id, updateFn)
}

// UpdateJobWorkerDaemon updates the job with the provided ID using the provided worker daemon update data.
// It modifies the job's status, result, and error message.
// The updated job is returned.
// functional options are used to modify the job's properties.
func (s *JobService) UpdateJobWorkerDaemon(ctx context.Context, id string, data ports.WorkerDaemonUpdateData) (ports.Job, error) {
	// Validate ID
	if len(strings.TrimSpace(id)) == 0 {
		return ports.Job{}, errors.New("job ID must be provided")
	}
	if _, err := uuid.Parse(id); err != nil {
		return ports.Job{}, errors.New("job ID must be a valid UUID")
	}

	updateFn := func(job *ports.Job) error {
		job.Status = data.Status

		if data.Status == ports.StatusCompleted && data.ErrorMessage != "" {
			return errors.New("completed status should not have an error message")
		}

		if data.Status == ports.StatusFailed && data.Result != "" {
			return errors.New("failed status should not have a result")
		}

		job.Result = data.Result
		job.ErrorMessage = data.ErrorMessage
		job.UpdatedAt = time.Now()
		return nil
	}

	return s.storage.UpdateJob(ctx, id, updateFn)
}

// isSimpleValidVersion checks if the image version string contains only valid characters.
func isSimpleValidVersion(version string) bool {
	for _, char := range version {
		if !(unicode.IsLetter(char) || unicode.IsDigit(char) || char == '.' || char == '-' || char == '_') {
			return false
		}
	}
	return true
}

// Helper function to validate if a given JobStatus is valid
func isValidStatus(status ports.JobStatus) bool {
	switch status {
	case ports.StatusQueued, ports.StatusScheduled, ports.StatusRunning, ports.StatusCompleted, ports.StatusFailed, ports.StatusCancelled:
		return true
	default:
		return false
	}
}
