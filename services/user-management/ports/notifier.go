package ports

import (
	"context"
)

type Notifier interface {
	UserManagementChanged(userManagement UserManagement, ctx context.Context)
}
