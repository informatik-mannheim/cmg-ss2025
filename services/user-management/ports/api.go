package ports

import (
	"context"
	"errors"
)

var ErrUserManagementNotFound = errors.New("userManagement not found")

type Api interface {
	Set(userManagement UserManagement, ctx context.Context) error
	Get(id string, ctx context.Context) (UserManagement, error)
}
