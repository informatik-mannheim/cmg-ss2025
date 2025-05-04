package core

import (
	"context"

	"github.com/informatik-mannheim/cmg-ss2025/services/job/ports"
)

type JobService struct {
	repo     ports.Repo
	notifier ports.Notifier
}

func NewJobService(repo ports.Repo, notifier ports.Notifier) *JobService {
	return &JobService{
		repo:     repo,
		notifier: notifier,
	}
}

var _ ports.Api = (*JobService)(nil)

func (s *JobService) Set(job ports.Job, ctx context.Context) error {
	err := s.repo.Store(job, ctx)
	if err != nil {
		return err
	}
	if s.notifier != nil {
		s.notifier.JobChanged(job, ctx)
	}
	return nil
}

func (s *JobService) Get(id string, ctx context.Context) (ports.Job, error) {
	job, err := s.repo.FindById(id, ctx)
	if err != nil {
		return ports.Job{}, err
	}
	if job.Id != id {
		return ports.Job{}, ports.ErrJobNotFound
	}
	return job, nil
}
