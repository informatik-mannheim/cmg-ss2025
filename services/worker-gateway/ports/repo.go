package ports

import (
	"context"
)

type Repo interface {
	Store(ctx context.Context) error
}
