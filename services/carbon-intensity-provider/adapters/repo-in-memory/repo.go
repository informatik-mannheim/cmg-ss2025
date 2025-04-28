package repo_in_memory

import (
	"context"

	"github.com/informatik-mannheim/cmg-ss2025/servicescarbon-intensity-provider/ports"
)

type Repo struct {
	carbonIntensityProviders map[string]ports.CarbonIntensityProvider
}

var _ ports.Repo = (*Repo)(nil)

func NewRepo() *Repo {
	return &Repo{
		carbonIntensityProviders: make(map[string]ports.CarbonIntensityProvider),
	}
}

func (r *Repo) Store(carbonIntensityProvider ports.CarbonIntensityProvider, ctx context.Context) error {
	r.carbonIntensityProviders[carbonIntensityProvider.Id] = carbonIntensityProvider
	return nil
}

func (r *Repo) FindById(id string, ctx context.Context) (ports.CarbonIntensityProvider, error) {
	carbonIntensityProvider, ok := r.carbonIntensityProviders[id]
	if !ok {
		return ports.CarbonIntensityProvider{}, ports.ErrCarbonIntensityProviderNotFound
	}
	return carbonIntensityProvider, nil
}
