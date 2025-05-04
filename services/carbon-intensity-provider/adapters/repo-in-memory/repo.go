package repo_in_memory

import (
	"context"

	"github.com/informatik-mannheim/cmg-ss2025/services/carbon-intensity-provider/model"
	"github.com/informatik-mannheim/cmg-ss2025/services/carbon-intensity-provider/ports"
)

type Repo struct {
	carbonIntensityProviders map[string]model.CarbonIntensityData
}

// Ensure Repo implements ports.Repo
var _ ports.Repo = (*Repo)(nil)

func NewRepo() *Repo {
	return &Repo{
		carbonIntensityProviders: make(map[string]model.CarbonIntensityData),
	}
}

func (r *Repo) Store(data model.CarbonIntensityData, ctx context.Context) error {
	r.carbonIntensityProviders[data.Zone] = data
	return nil
}

func (r *Repo) FindById(id string, ctx context.Context) (model.CarbonIntensityData, error) {
	data, ok := r.carbonIntensityProviders[id]
	if !ok {
		return model.CarbonIntensityData{}, ports.ErrCarbonIntensityProviderNotFound
	}
	return data, nil
}

func (r *Repo) FindAll(ctx context.Context) ([]model.CarbonIntensityData, error) {
	result := make([]model.CarbonIntensityData, 0, len(r.carbonIntensityProviders))
	for _, data := range r.carbonIntensityProviders {
		result = append(result, data)
	}
	return result, nil
}
