package repo_in_memory

import (
	"context"

	"github.com/informatik-mannheim/cmg-ss2025/services/carbon-intensity-provider/ports"
)

type Repo struct {
	carbonIntensityProviders map[string]ports.CarbonIntensityData
}

// Ensure Repo implements ports.Repo
var _ ports.Repo = (*Repo)(nil)

func NewRepo() *Repo {
	return &Repo{
		carbonIntensityProviders: make(map[string]ports.CarbonIntensityData),
	}
}

func (r *Repo) Store(data ports.CarbonIntensityData, ctx context.Context) error {
	r.carbonIntensityProviders[data.Zone] = data
	return nil
}

func (r *Repo) FindById(id string, ctx context.Context) (ports.CarbonIntensityData, error) {
	data, ok := r.carbonIntensityProviders[id]
	if !ok {
		return ports.CarbonIntensityData{}, ports.ErrCarbonIntensityProviderNotFound
	}
	return data, nil
}

func (r *Repo) FindAll(ctx context.Context) ([]ports.CarbonIntensityData, error) {
	result := make([]ports.CarbonIntensityData, 0, len(r.carbonIntensityProviders))
	for _, data := range r.carbonIntensityProviders {
		result = append(result, data)
	}
	return result, nil
}
