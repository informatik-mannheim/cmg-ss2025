package ports

import (
	"context"
)

type Api interface {
	Set(carbonIntensityProvider CarbonIntensityProvider, ctx context.Context) error
	Get(id string, ctx context.Context) (CarbonIntensityProvider, error)
}
