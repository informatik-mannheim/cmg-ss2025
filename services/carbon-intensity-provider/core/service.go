package core

import (
	"context"

	"github.com/informatik-mannheim/cmg-ss2025/services/carbon-intensity-provider/ports"
)

type CarbonIntensityProvider struct {
	repo     ports.Repo
	notifier ports.Notifier
}

func NewCarbonIntensityProvider(repo ports.Repo, notifier ports.Notifier) *CarbonIntensityProvider {
	return &CarbonIntensityProvider{
		repo:     repo,
		notifier: notifier,
	}
}

var _ ports.Api = (*CarbonIntensityProvider)(nil)

func (s *CarbonIntensityProvider) Set(carbonIntensityProvider ports.CarbonIntensityProvider, ctx context.Context) error {
	err := s.repo.Store(carbonIntensityProvider, ctx)
	if err != nil {
		return err
	}
	if s.notifier != nil {
		s.notifier.CarbonIntensityProviderChanged(carbonIntensityProvider, ctx)
	}
	return nil
}

func (s *CarbonIntensityProvider) Get(id string, ctx context.Context) (ports.CarbonIntensityProvider, error) {
	carbonIntensityProvider, err := s.repo.FindById(id, ctx)
	if err != nil {
		return ports.CarbonIntensityProvider{}, err
	}
	if carbonIntensityProvider.Id != id {
		return ports.CarbonIntensityProvider{}, ports.ErrCarbonIntensityProviderNotFound
	}
	return carbonIntensityProvider, nil
}
