package core

import (
	"context"

	"github.com/informatik-mannheim/cmg-ss2025/services/user-management/ports"
)

type UserManagementService struct {
	repo     ports.Repo
	notifier ports.Notifier
}

func NewUserManagementService(repo ports.Repo, notifier ports.Notifier) *UserManagementService {
	return &UserManagementService{
		repo:     repo,
		notifier: notifier,
	}
}

var _ ports.Api = (*UserManagementService)(nil)

func (s *UserManagementService) Set(userManagement ports.UserManagement, ctx context.Context) error {
	err := s.repo.Store(userManagement, ctx)
	if err != nil {
		return err
	}
	if s.notifier != nil {
		s.notifier.UserManagementChanged(userManagement, ctx)
	}
	return nil
}

func (s *UserManagementService) Get(id string, ctx context.Context) (ports.UserManagement, error) {
	userManagement, err := s.repo.FindById(id, ctx)
	if err != nil {
		return ports.UserManagement{}, err
	}
	if userManagement.Id != id {
		return ports.UserManagement{}, ports.ErrUserManagementNotFound
	}
	return userManagement, nil
}
