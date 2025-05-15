package core

import (
	"context"
	"log"

	"github.com/informatik-mannheim/cmg-ss2025/services/worker-gateway/ports"
)

type WorkerGatewayService struct {
	notifier ports.Notifier
}

func NewWorkerGatewayService(notifier ports.Notifier) *WorkerGatewayService {
	return &WorkerGatewayService{
		notifier: notifier,
	}
}

func (s *WorkerGatewayService) Heartbeat(ctx context.Context, req ports.HeartbeatRequest) ([]ports.Job, error) {
	log.Printf("Heartbeat received: %s is %s", req.WorkerID, req.Status)

	if err := s.notifier.UpdateWorkerStatus(ctx, req); err != nil {
		log.Printf("update worker status error: %v", err)
	}

	// Get scheduled jobs if available
	if req.Status == "AVAILABLE" {
		jobs, err := s.notifier.FetchScheduledJobs(ctx)
		if err != nil {
			log.Printf("error getting jobs: %v", err)
			return nil, nil // Gateway still accepts heartbeat
		}

		// Filter for this worker
		var assigned []ports.Job
		for _, job := range jobs {
			if job.WorkerID == req.WorkerID {
				assigned = append(assigned, job)
			}
		}
		return assigned, nil
	}

	// Comuting
	return nil, nil
}

func (s *WorkerGatewayService) SubmitResult(ctx context.Context, result ports.ResultRequest) error {
	log.Printf("Result for job %s received", result.JobID)

	return s.notifier.UpdateJob(ctx, result)
}

func (s *WorkerGatewayService) RegisterWorker(ctx context.Context, req ports.RegisterRequest) error {
	log.Printf("Registering worker: %s", req.ID)

	return s.notifier.RegisterWorker(ctx, req)
}
