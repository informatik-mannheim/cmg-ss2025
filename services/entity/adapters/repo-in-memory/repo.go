package repo_in_memory

import (
	"context"

	"gitty.informatik.hs-mannheim.de/steger/cmg-ws202425/services/entity/ports"
)

type Repo struct {
	entities map[string]ports.Entity
}

var _ ports.Repo = (*Repo)(nil)

func NewRepo() *Repo {
	return &Repo{
		entities: make(map[string]ports.Entity),
	}
}

func (r *Repo) Store(entity ports.Entity, ctx context.Context) error {
	r.entities[entity.Id] = entity
	return nil
}

func (r *Repo) FindById(id string, ctx context.Context) (ports.Entity, error) {
	entity, ok := r.entities[id]
	if !ok {
		return ports.Entity{}, ports.ErrEntityNotFound
	}
	return entity, nil
}
