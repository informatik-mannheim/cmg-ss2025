package core

import (
	"context"

	"github.com/informatik-mannheim/cmg-ss2025/pkg/logging"
	"github.com/informatik-mannheim/cmg-ss2025/services/worker-gateway/ports"
)

type WorkerGatewayService struct {
	registry ports.RegistryService
	job      ports.JobService
	user     ports.UserService
}

func NewWorkerGatewayService(registry ports.RegistryService, job ports.JobService, user ports.UserService) *WorkerGatewayService {
	return &WorkerGatewayService{registry: registry, job: job, user: user}
}

func (s *WorkerGatewayService) Heartbeat(ctx context.Context, req ports.HeartbeatRequest, token string) ([]ports.Job, error) {

	logging.From(ctx).Debug("Heartbeat received", "workerID", req.WorkerID, "status", req.Status)

	if err := s.registry.UpdateWorkerStatus(ctx, req, token); err != nil {
		logging.From(ctx).Error("UpdateWorkerStatus failed", "error", err)
		return nil, err
	}

	if req.Status == "AVAILABLE" {
		jobs, err := s.job.FetchScheduledJobs(ctx, token)
		if err != nil {
			logging.From(ctx).Error("Error fetching jobs", "error", err)
			return nil, err
		}
		logging.From(ctx).Debug("Provided jobs", "jobs", jobs)

		var filteredJobs []ports.Job
		for _, job := range jobs {
			if job.WorkerID == req.WorkerID {
				filteredJobs = append(filteredJobs, job)
			}
		}

		logging.From(ctx).Debug("Filtered jobs", "filteredJobs", filteredJobs)

		return filteredJobs, nil
	}

	return nil, nil
}

func (s *WorkerGatewayService) Result(ctx context.Context, result ports.ResultRequest, token string) error {

	logging.From(ctx).Debug("Result received", "jobID", result.JobID)
	return s.job.UpdateJob(ctx, result, token)
}

func (s *WorkerGatewayService) Register(ctx context.Context, req ports.RegisterRequest) (*ports.RegisterRespose, error) {
	logging.From(ctx).Debug("Registering worker", "zone", req.Zone)

	tokenReq := ports.GetTokenRequest{
		Secret: req.Key,
	}

	tokenResp, err := s.user.GetToken(ctx, tokenReq)
	if err != nil {
		logging.From(ctx).Error("Failed to register provider with user service", "error", err)
		return nil, err
	}

	logging.From(ctx).Debug("Received secret from user service")

	regResp, err := s.registry.RegisterWorker(ctx, req, tokenResp.Token)
	if err != nil {
		logging.From(ctx).Error("Worker registration failed", "error", err)
		return nil, err
	}

	logging.From(ctx).Debug("Worker registered", "workerID", regResp.ID, "zone", regResp.Zone)
	return regResp, nil
}
