package core

import (
	"context"

	uuid "github.com/google/uuid"
	"github.com/informatik-mannheim/cmg-ss2025/services/worker-registry/ports"
)

type WorkerRegistryService struct {
	repo          ports.Repo
	notifier      ports.Notifier
	zoneValidator ports.ZoneValidator
}

func NewWorkerRegistryService(repo ports.Repo, notifier ports.Notifier, validator ports.ZoneValidator) *WorkerRegistryService {
	return &WorkerRegistryService{
		repo:          repo,
		notifier:      notifier,
		zoneValidator: validator,
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
	if !s.zoneValidator.IsValidZone(zone, ctx) {
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

	s.notifier.WorkerCreated(newWorker, ctx)
	return newWorker, nil
}

func (s *WorkerRegistryService) UpdateWorkerStatus(id string, status ports.WorkerStatus, ctx context.Context) (ports.Worker, error) {
	newWorker, err := s.repo.UpdateWorkerStatus(id, status, ctx)
	if err != nil {
		return ports.Worker{}, err
	}
	s.notifier.WorkerStatusChanged(newWorker, ctx)
	return newWorker, nil
}
