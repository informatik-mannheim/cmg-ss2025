package ports

import (
	"context"

	"github.com/informatik-mannheim/cmg-ss2025/services/carbon-intensity-provider/model"
)

type Repo interface {
	Store(data model.CarbonIntensityData, ctx context.Context) error
	FindById(id string, ctx context.Context) (model.CarbonIntensityData, error)
	FindAll(ctx context.Context) ([]model.CarbonIntensityData, error)
}
