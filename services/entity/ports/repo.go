package ports

import (
	"context"
)

type Repo interface {
	Store(entity Entity, ctx context.Context) error
	FindById(id string, ctx context.Context) (Entity, error)
}
