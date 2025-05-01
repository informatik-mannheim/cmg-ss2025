package core

import (
	"context"

	"github.com/informatik-mannheim/cmg-ss2025/services/consumer-gateway/ports"
)

type ConsumerService struct {
	repo     ports.Repo
	notifier ports.Notifier
}

func NewConsumerService(repo ports.Repo, notifier ports.Notifier) *ConsumerService {
	return &ConsumerService{
		repo:     repo,
		notifier: notifier,
	}
}

var _ ports.Api = (*ConsumerService)(nil)

func (s *ConsumerService) Set(consumer ports.Consumer, ctx context.Context) error {
	err := s.repo.Store(consumer, ctx)
	if err != nil {
		return err
	}
	if s.notifier != nil {
		s.notifier.ConsumerChanged(consumer, ctx)
	}
	return nil
}

func (s *ConsumerService) Get(id string, ctx context.Context) (ports.Consumer, error) {
	consumer, err := s.repo.FindById(id, ctx)
	if err != nil {
		return ports.Consumer{}, err
	}
	if consumer.Id != id {
		return ports.Consumer{}, ports.ErrConsumerNotFound
	}
	return consumer, nil
}
