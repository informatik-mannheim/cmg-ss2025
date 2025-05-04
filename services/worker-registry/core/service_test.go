package core

import (
	"context"
	"testing"

	notifier "github.com/informatik-mannheim/cmg-ss2025/services/worker-registry/adapters/notifier"
	repo_in_memory "github.com/informatik-mannheim/cmg-ss2025/services/worker-registry/adapters/repo-in-memory"
	"github.com/informatik-mannheim/cmg-ss2025/services/worker-registry/ports"
)

func TestGetWorkers(t *testing.T) {
	repo := repo_in_memory.NewRepo()
	notifier := notifier.NewHttpNotifier()
	service := NewWorkerRegistryService(repo, notifier)

	service.CreateWorker("DE", context.Background())
	service.CreateWorker("EN", context.Background())

	tests := []struct {
		name          string
		status        string
		zone          string
		expectedCount int
	}{
		{"all workers", "", "", 2},
		{"by status status 'AVAILABLE' only", "AVAILABLE", "", 2},
		{"by status 'RUNNING' only", "RUNNING", "", 0},
		{"by zone 'DE' only", "", "DE", 1},
		{"by zone 'EN' only", "", "EN", 1},
		{"by status and zone", "AVAILABLE", "EN", 1},
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
	notifier := notifier.NewHttpNotifier()
	service := NewWorkerRegistryService(repo, notifier)

	worker, _ := service.CreateWorker("DE", context.Background())

	t.Run("existing worker", func(t *testing.T) {
		result, err := service.GetWorkerById(worker.Id, context.Background())
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result.Id != worker.Id {
			t.Errorf("expected ID %s, got %s", worker.Id, result.Id)
		}
	})

	t.Run("non-existent worker", func(t *testing.T) {
		_, err := service.GetWorkerById("9999", context.Background())
		if err != ports.ErrWorkerNotFound {
			t.Errorf("expected ErrWorkerNotFound, got %v", err)
		}
	})
}

func TestUpdateWorkerStatus(t *testing.T) {
	repo := repo_in_memory.NewRepo()
	notifier := notifier.NewHttpNotifier()
	service := NewWorkerRegistryService(repo, notifier)

	worker, _ := service.CreateWorker("ZoneX", context.Background())

	t.Run("valid status update", func(t *testing.T) {
		updated, err := service.UpdateWorkerStatus(worker.Id, "RUNNING", context.Background())
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if updated.Status != "RUNNING" {
			t.Errorf("expected status RUNNING, got %s", updated.Status)
		}
	})

	t.Run("invalid status update", func(t *testing.T) {
		_, err := service.UpdateWorkerStatus(worker.Id, "INVALID_STATUS", context.Background())
		if err != ports.ErrUpdatingWorkerFailed {
			t.Errorf("expected ErrUpdatingWorkerFailed, got %v", err)
		}
	})

	t.Run("non-existent worker", func(t *testing.T) {
		_, err := service.UpdateWorkerStatus("9999", "AVAILABLE", context.Background())
		if err != ports.ErrWorkerNotFound {
			t.Errorf("expected ErrWorkerNotFound, got %v", err)
		}
	})
}
