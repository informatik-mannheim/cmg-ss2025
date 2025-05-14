package repo_in_memory

import (
	"context"
	"errors"

	"github.com/informatik-mannheim/cmg-ss2025/services/job/ports"
	"github.com/informatik-mannheim/cmg-ss2025/services/job/utils"
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
		if len(status) == 0 || utils.ContainsStatus(status, job.Status) {
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
