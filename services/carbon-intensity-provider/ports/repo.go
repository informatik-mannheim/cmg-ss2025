package ports

import (
	"context"
<<<<<<< HEAD
	"errors"

	"github.com/informatik-mannheim/cmg-ss2025/services/carbon-intensity-provider/model"
)

var ErrCarbonIntensityProviderNotFound = errors.New("carbon intensity provider not found")

type Repo interface {
	Store(data model.CarbonIntensityData, ctx context.Context) error
	FindById(id string, ctx context.Context) (model.CarbonIntensityData, error)
	FindAll(ctx context.Context) ([]model.CarbonIntensityData, error)
=======
)

type Repo interface {
	Store(carbonIntensityProvider CarbonIntensityProvider, ctx context.Context) error
	FindById(id string, ctx context.Context) (CarbonIntensityProvider, error)
>>>>>>> origin
}
