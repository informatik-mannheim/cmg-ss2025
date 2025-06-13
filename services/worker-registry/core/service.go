package core

import (
	"context"
	"fmt"

	uuid "github.com/google/uuid"
	"github.com/informatik-mannheim/cmg-ss2025/pkg/logging"
	"github.com/informatik-mannheim/cmg-ss2025/services/worker-registry/ports"
)

type WorkerRegistryService struct {
	repo       ports.Repo
	zoneClient ports.ZoneClient
}

func NewWorkerRegistryService(repo ports.Repo, zoneClient ports.ZoneClient) *WorkerRegistryService {
	return &WorkerRegistryService{
		repo:       repo,
		zoneClient: zoneClient,
	}
}

var _ ports.Api = (*WorkerRegistryService)(nil)

func (s *WorkerRegistryService) GetWorkers(status ports.WorkerStatus, zone string, ctx context.Context) ([]ports.Worker, error) {
	return s.repo.GetWorkers(status, zone, ctx)
}

func (s *WorkerRegistryService) GetWorkerById(id string, ctx context.Context) (ports.Worker, error) {
	return s.repo.GetWorkerById(id, ctx)
}

func (s *WorkerRegistryService) CreateWorker(zone string, ctx context.Context) (ports.Worker, error) {
	if !s.zoneClient.IsValidZone(zone, ctx) {
		return ports.Worker{}, ports.NewErrCreatingWorkerFailedInvalidZone(zone)
	}

	newWorker := ports.Worker{
		Id:     uuid.NewString(),
		Status: ports.StatusAvailable,
		Zone:   zone,
	}
	err := s.repo.CreateWorker(newWorker, ctx)
	if err != nil {
		return ports.Worker{}, err
	}

	resultMessage := fmt.Sprintf("New worker created: ID=%s, STATUS=%s, ZONE=%s", newWorker.Id, newWorker.Status, newWorker.Zone)
	logging.Debug(resultMessage)
	return newWorker, nil
}

func (s *WorkerRegistryService) UpdateWorkerStatus(id string, status ports.WorkerStatus, ctx context.Context) (ports.Worker, error) {
	newWorker, err := s.repo.UpdateWorkerStatus(id, status, ctx)
	if err != nil {
		return ports.Worker{}, err
	}

	resultMessage := fmt.Sprintf("[Notifier] Changed status from Worker with ID '%s' to status '%s'.", newWorker.Id, newWorker.Status)
	logging.Debug(resultMessage)
	return newWorker, nil
}
