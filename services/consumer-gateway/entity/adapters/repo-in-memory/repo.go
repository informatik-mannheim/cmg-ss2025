package repo_in_memory

import (
	"context"

	"github.com/informatik-mannheim/cmg-ss2025/services/consumer-gateway/ports"
)

type Repo struct {
	consumers map[string]ports.Consumer
}

var _ ports.Repo = (*Repo)(nil)

func NewRepo() *Repo {
	return &Repo{
		consumers: make(map[string]ports.Consumer),
	}
}

func (r *Repo) Store(entity ports.Consumer, ctx context.Context) error {
	r.consumers[Consumer.Id] = entity
	return nil
}

func (r *Repo) FindById(id string, ctx context.Context) (ports.Entity, error) {
	consumer, ok := r.consumers[id]
	if !ok {
		return ports.Consumer{}, ports.ErrConsumerNotFound
	}
	return consumer, nil
}
