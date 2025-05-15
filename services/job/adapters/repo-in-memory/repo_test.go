package repo_in_memory_test

import (
	"context"
	"testing"

	repo_in_memory "github.com/informatik-mannheim/cmg-ss2025/services/job/adapters/repo-in-memory"
	"github.com/informatik-mannheim/cmg-ss2025/services/job/ports"
)

func TestCreateJob(t *testing.T) {
	storage := repo_in_memory.NewMockJobStorage()
	job := ports.Job{Id: "1", JobName: "TestJob", Status: ports.StatusQueued}

	// Test creating a new job
	err := storage.CreateJob(context.Background(), job)
	if err != nil {
		t.Fatalf("expected no error but got %v", err)
	}
}

func TestGetJob_ExistingJob(t *testing.T) {
	storage := repo_in_memory.NewMockJobStorage()
	job := ports.Job{Id: "1", JobName: "TestJob", Status: ports.StatusQueued}
	_ = storage.CreateJob(context.Background(), job)

	// Test retrieving the existing job
	retrievedJob, err := storage.GetJob(context.Background(), "1")
	if err != nil {
		t.Fatalf("expected no error but got %v", err)
	}
	if retrievedJob.Id != job.Id {
		t.Errorf("expected %v, got %v", job.Id, retrievedJob.Id)
	}
}

func TestGetJob_NonExistingJob(t *testing.T) {
	storage := repo_in_memory.NewMockJobStorage()

	// Test retrieving a non-existing job
	_, err := storage.GetJob(context.Background(), "non-existing-id")
	if err != ports.ErrJobNotFound {
		t.Fatalf("expected %v error but got %v", ports.ErrJobNotFound, err)
	}
}

func TestGetJobs_NoFilter(t *testing.T) {
	storage := repo_in_memory.NewMockJobStorage()
	job1 := ports.Job{Id: "1", JobName: "TestJob1", Status: ports.StatusQueued}
	job2 := ports.Job{Id: "2", JobName: "TestJob2", Status: ports.StatusCompleted}
	_ = storage.CreateJob(context.Background(), job1)
	_ = storage.CreateJob(context.Background(), job2)

	// Test retrieving all jobs with no filter
	allJobs, err := storage.GetJobs(context.Background(), nil)
	if err != nil {
		t.Fatalf("expected no error but got %v", err)
	}
	if len(allJobs) != 2 {
		t.Errorf("expected 2 jobs, got %d", len(allJobs))
	}
}

func TestGetJobs_WithFilter(t *testing.T) {
	storage := repo_in_memory.NewMockJobStorage()
	job1 := ports.Job{Id: "1", JobName: "TestJob1", Status: ports.StatusQueued}
	job2 := ports.Job{Id: "2", JobName: "TestJob2", Status: ports.StatusCompleted}
	_ = storage.CreateJob(context.Background(), job1)
	_ = storage.CreateJob(context.Background(), job2)

	// Test retrieving jobs with a status filter
	queuedJobs, err := storage.GetJobs(context.Background(), []ports.JobStatus{ports.StatusQueued})
	if err != nil {
		t.Fatalf("expected no error but got %v", err)
	}
	if len(queuedJobs) != 1 {
		t.Errorf("expected 1 queued job, got %d", len(queuedJobs))
	}
}

func TestUpdateJob_ExistingJob(t *testing.T) {
	storage := repo_in_memory.NewMockJobStorage()
	job := ports.Job{Id: "1", JobName: "TestJob", Status: ports.StatusQueued}
	_ = storage.CreateJob(context.Background(), job)

	updatedJob := ports.Job{Id: "1", JobName: "UpdatedTestJob", Status: ports.StatusRunning}
	result, err := storage.UpdateJob(context.Background(), "1", updatedJob)
	if err != nil {
		t.Fatalf("expected no error but got %v", err)
	}
	if result.JobName != "UpdatedTestJob" {
		t.Errorf("expected job name %v, got %v", "UpdatedTestJob", result.JobName)
	}
}

func TestUpdateJob_NonExistingJob(t *testing.T) {
	storage := repo_in_memory.NewMockJobStorage()
	updatedJob := ports.Job{Id: "1", JobName: "UpdatedTestJob", Status: ports.StatusRunning}

	// Test updating a non-existing job
	_, err := storage.UpdateJob(context.Background(), "non-existing-id", updatedJob)
	if err != ports.ErrJobNotFound {
		t.Fatalf("expected %v error but got %v", ports.ErrJobNotFound, err)
	}
}
