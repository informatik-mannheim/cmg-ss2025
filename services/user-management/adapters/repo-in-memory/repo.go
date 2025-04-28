package repo_in_memory

import (
	"context"

	"github.com/informatik-mannheim/cmg-ss2025/services/user-management/ports"
)

type Repo struct {
	userManagements map[string]ports.UserManagement
}

var _ ports.Repo = (*Repo)(nil)

func NewRepo() *Repo {
	return &Repo{
		userManagements: make(map[string]ports.UserManagement),
	}
}

func (r *Repo) Store(userManagement ports.UserManagement, ctx context.Context) error {
	r.userManagements[userManagement.Id] = userManagement
	return nil
}

func (r *Repo) FindById(id string, ctx context.Context) (ports.UserManagement, error) {
	userManagement, ok := r.userManagements[id]
	if !ok {
		return ports.UserManagement{}, ports.ErrUserManagementNotFound
	}
	return userManagement, nil
}
