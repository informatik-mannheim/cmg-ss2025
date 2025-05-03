package ports

import (
	"context"
)

type Repo interface {
	Store(carbonIntensityProvider CarbonIntensityProvider, ctx context.Context) error
	FindById(id string, ctx context.Context) (CarbonIntensityProvider, error)
}
