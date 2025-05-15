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
	var _ ports.JobService = (*JobService)(nil)

	return &JobService{
		storage: storage,
	}, nil
}

// GetJobs retrieves jobs based on their status.
// It returns a slice of jobs that match the provided status.
// If no status is provided, it returns all jobs.
func (s *JobService) GetJobs(ctx context.Context, status []ports.JobStatus) ([]ports.Job, error) {

	if len(status) == 0 {
		return s.storage.GetJobs(ctx, nil)
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
		return ports.Job{}, ports.ErrNotExistingJobName
	}
	if strings.TrimSpace(jobCreate.CreationZone) == "" {
		return ports.Job{}, ports.ErrNotExistingZone
	}
	if strings.TrimSpace(jobCreate.Image.Name) == "" {
		return ports.Job{}, ports.ErrNotExistingImageName
	}
	if !isSimpleValidVersion(jobCreate.Image.Version) {
		return ports.Job{}, ports.ErrImageVersionIsInvalid
	}

	for key, value := range jobCreate.Parameters {
		if strings.TrimSpace(key) == "" || strings.TrimSpace(value) == "" {
			return ports.Job{}, ports.ErrParamKeyValueEmpty
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
		return ports.Job{}, ports.ErrNotExistingID
	}
	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return ports.Job{}, ports.ErrInvalidIDFormat
	}
	// Check if the job exists in the storage
	if _, err := s.storage.GetJob(ctx, id); err != nil {
		return ports.Job{}, ports.ErrJobNotFound
	}

	return s.storage.GetJob(ctx, id)
}

// GetJobOutcome retrieves the outcome of a job by its ID.
// It returns the job name, status, result, error message, compute zone, carbon intensity, and carbon savings.
// If the job is not found, it returns an error.
// This method is useful for getting detailed information about a job's execution.
func (s *JobService) GetJobOutcome(ctx context.Context, id string) (ports.JobOutcome, error) {
	var job ports.Job

	// Check for empty or whitespace-only ID
	if len(strings.TrimSpace(id)) == 0 {
		return ports.JobOutcome{}, ports.ErrNotExistingID
	}
	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return ports.JobOutcome{}, ports.ErrInvalidIDFormat
	}
	// Check if the job exists in the storage
	job, err := s.storage.GetJob(ctx, id)
	if err != nil {
		return ports.JobOutcome{}, ports.ErrJobNotFound
	}

	return ports.JobOutcome{
		JobName:         job.JobName,
		Status:          job.Status,
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
		return ports.Job{}, ports.ErrNotExistingID
	}
	if _, err := uuid.Parse(id); err != nil {
		return ports.Job{}, ports.ErrInvalidIDFormat
	}
	if strings.TrimSpace(string(data.Status)) == "" {
		return ports.Job{}, ports.ErrNotExistingStatus
	}

	// Validate WorkerID
	if len(strings.TrimSpace(data.WorkerID)) == 0 {
		return ports.Job{}, ports.ErrNotExistingWorkerID
	}

	// Validate CarbonIntensity
	if data.CarbonIntensity < 0 {
		return ports.Job{}, ports.ErrCarbonIsNegative
	}

	updated_job, err := s.GetJob(ctx, id)

	if err != nil {
		return ports.Job{}, err
	}
	updated_job.WorkerID = data.WorkerID
	updated_job.ComputeZone = data.ComputeZone
	updated_job.CarbonIntensity = data.CarbonIntensity
	updated_job.CarbonSaving = data.CarbonSaving
	updated_job.Status = data.Status
	updated_job.UpdatedAt = time.Now()

	return s.storage.UpdateJob(ctx, id, updated_job)
}

// UpdateJobWorkerDaemon updates the job with the provided ID using the provided worker daemon update data.
// It modifies the job's status, result, and error message.
// The updated job is returned.
// functional options are used to modify the job's properties.
func (s *JobService) UpdateJobWorkerDaemon(ctx context.Context, id string, data ports.WorkerDaemonUpdateData) (ports.Job, error) {
	// Validate ID
	if len(strings.TrimSpace(id)) == 0 {
		return ports.Job{}, ports.ErrNotExistingID
	}
	if _, err := uuid.Parse(id); err != nil {
		return ports.Job{}, ports.ErrInvalidIDFormat
	}
	if strings.TrimSpace(string(data.Status)) == "" {
		return ports.Job{}, ports.ErrNotExistingStatus
	}
	if data.Status == ports.StatusFailed && strings.TrimSpace(data.ErrorMessage) == "" {
		return ports.Job{}, ports.ErrErrorMessageEmpty
	}

	updated_job, err := s.GetJob(ctx, id)

	if err != nil {
		return ports.Job{}, err
	}
	updated_job.Status = data.Status
	updated_job.Result = data.Result
	updated_job.ErrorMessage = data.ErrorMessage
	updated_job.UpdatedAt = time.Now()

	return s.storage.UpdateJob(ctx, id, updated_job)
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
