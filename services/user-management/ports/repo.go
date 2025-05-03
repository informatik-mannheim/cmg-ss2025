package ports

import (
	"context"
)

type Repo interface {
	Store(userManagement UserManagement, ctx context.Context) error
	FindById(id string, ctx context.Context) (UserManagement, error)
}
