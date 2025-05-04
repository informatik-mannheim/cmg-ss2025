package ports

import (
	"context"
)

type Repo interface {
	Store(job Job, ctx context.Context) error
	FindById(id string, ctx context.Context) (Job, error)
}
