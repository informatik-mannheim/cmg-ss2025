package repo_in_memory

import (
	"context"

	"github.com/informatik-mannheim/cmg-ss2025/services/worker-gateway/ports"
)

type Repo struct {
}

var _ ports.Repo = (*Repo)(nil)

func NewRepo() *Repo {
	return nil
}

func (r *Repo) Store(ectx context.Context) error {
	return nil
}
