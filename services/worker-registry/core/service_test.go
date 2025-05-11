package core

import (
	"context"
	"fmt"
	"testing"

	notifier "github.com/informatik-mannheim/cmg-ss2025/services/worker-registry/adapters/notifier"
	repo_in_memory "github.com/informatik-mannheim/cmg-ss2025/services/worker-registry/adapters/repo-in-memory"
	validator "github.com/informatik-mannheim/cmg-ss2025/services/worker-registry/adapters/zone-validator"
	"github.com/informatik-mannheim/cmg-ss2025/services/worker-registry/ports"
)

func TestCreateWorker(t *testing.T) {
	repo := repo_in_memory.NewRepo()
	notifier := notifier.NewNotifier()
	zoneValidator := validator.NewZoneValidator()
	service := NewWorkerRegistryService(repo, notifier, zoneValidator)

	t.Run("create worker with valid zone", func(t *testing.T) {
		worker, err := service.CreateWorker("EN", context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if worker.Zone != "EN" {
			t.Errorf("expected zone 'EN', got %v", worker.Zone)
		}
		if worker.Status != ports.StatusAvailable {
			t.Errorf("expected status 'AVAILABLE', got %v", worker.Status)
		}
		if worker.Id == "" {
			t.Errorf("expected a non-empty ID")
		}
	})

	t.Run("create worker with invalid zone", func(t *testing.T) {
		worker, err := service.CreateWorker("CMG", context.Background())
		if err == nil {
			t.Fatalf("expected error, got worker with ID %v", worker.Id)
		}
		expectedError := "creating worker failed due to invalid 'zone' CMG"
		if err.Error() != expectedError {
			t.Errorf("expected error: %v, got: %v", expectedError, err)
		}
	})
}

func TestGetWorkers(t *testing.T) {
	repo := repo_in_memory.NewRepo()
	notifier := notifier.NewNotifier()
	zoneValidator := validator.NewZoneValidator()
	service := NewWorkerRegistryService(repo, notifier, zoneValidator)

	service.CreateWorker("DE", context.Background())
	service.CreateWorker("EN", context.Background())

	tests := []struct {
		name          string
		status        ports.WorkerStatus
		zone          string
		expectedCount int
	}{
		{"all workers", "", "", 2},
		{"by status status 'AVAILABLE' only", ports.StatusAvailable, "", 2},
		{"by status 'RUNNING' only", ports.StatusRunning, "", 0},
		{"by zone 'DE' only", ports.StatusAvailable, "DE", 1},
		{"by zone 'EN' only", ports.StatusAvailable, "EN", 1},
		{"by status and zone", ports.StatusAvailable, "EN", 1},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			workers, err := service.GetWorkers(test.status, test.zone, context.Background())
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if len(workers) != test.expectedCount {
				t.Errorf("expected %d workers, got %d", test.expectedCount, len(workers))
			}
		})
	}
}

func TestGetWorkerById(t *testing.T) {
	repo := repo_in_memory.NewRepo()
	notifier := notifier.NewNotifier()
	zoneValidator := validator.NewZoneValidator()
	service := NewWorkerRegistryService(repo, notifier, zoneValidator)

	worker, _ := service.CreateWorker("DE", context.Background())

	t.Run("existing worker", func(t *testing.T) {
		result, err := service.GetWorkerById(worker.Id, context.Background())
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result.Id != worker.Id {
			t.Errorf("expected ID %v, got %v", worker.Id, result.Id)
		}
	})

	t.Run("non-existent worker", func(t *testing.T) {
		_, err := service.GetWorkerById("9999", context.Background())
		expectedError := "Worker with ID 9999 not found"
		if err == nil || err.Error() != expectedError {
			t.Errorf("expected error: %v, got: %v", expectedError, err)
		}
	})
}

func TestUpdateWorkerStatus(t *testing.T) {
	repo := repo_in_memory.NewRepo()
	notifier := notifier.NewNotifier()
	zoneValidator := validator.NewZoneValidator()
	service := NewWorkerRegistryService(repo, notifier, zoneValidator)

	worker, _ := service.CreateWorker("DE", context.Background())

	t.Run("valid status update", func(t *testing.T) {
		updated, err := service.UpdateWorkerStatus(worker.Id, "RUNNING", context.Background())
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if updated.Status != "RUNNING" {
			t.Errorf("expected status RUNNING, got %v", updated.Status)
		}
	})

	t.Run("invalid status update", func(t *testing.T) {
		_, err := service.UpdateWorkerStatus(worker.Id, "INVALID_STATUS", context.Background())
		expectedError := fmt.Sprintf("invalid status ('AVAILABLE' or 'RUNNING') for worker with ID %s", worker.Id)
		if err == nil || err.Error() != expectedError {
			t.Errorf("expected error: %v, got: %v", expectedError, err)
		}
	})

	t.Run("non-existent worker", func(t *testing.T) {
		_, err := service.UpdateWorkerStatus("9999", "AVAILABLE", context.Background())
		expectedError := "Worker with ID 9999 not found"
		if err == nil || err.Error() != expectedError {
			t.Errorf("expected error: %v, got: %v", expectedError, err)
		}
	})
}
