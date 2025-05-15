package repo_in_memory

import (
	"context"

	"github.com/informatik-mannheim/cmg-ss2025/services/job/ports"
	"github.com/informatik-mannheim/cmg-ss2025/services/job/utils"
)

// MockJobStorage is a mock implementation of the JobStorage interface
type MockJobStorage struct {
	jobs map[string]ports.Job
}

func NewMockJobStorage() *MockJobStorage {
	var _ ports.JobStorage = (*MockJobStorage)(nil)

	return &MockJobStorage{
		jobs: make(map[string]ports.Job),
	}
}

func (m *MockJobStorage) GetJobs(ctx context.Context, status []ports.JobStatus) ([]ports.Job, error) {
	var results []ports.Job

	if len(status) == 0 {
		// Return all jobs if no specific status is provided
		for _, job := range m.jobs {
			results = append(results, job)
		}
		return results, nil
	}

	for _, job := range m.jobs {
		if utils.ContainsStatus(status, job.Status) {
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
		return ports.Job{}, ports.ErrJobNotFound
	}
	return job, nil
}

func (m *MockJobStorage) UpdateJob(ctx context.Context, id string, updatedJob ports.Job) (ports.Job, error) {
	_, exists := m.jobs[id]
	if !exists {
		return ports.Job{}, ports.ErrJobNotFound
	}
	m.jobs[id] = updatedJob
	return updatedJob, nil
}
