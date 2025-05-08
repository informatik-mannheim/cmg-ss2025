package ports

import (
	"context"
)

type Repo interface {
	Store(consumer Consumer, ctx context.Context) error
	FindById(id string, ctx context.Context) (Consumer, error)
}
