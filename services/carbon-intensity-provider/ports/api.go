package ports

import (
	"context"
	"errors"
)

var ErrCarbonIntensityProviderNotFound = errors.New("carbonIntensityProvider not found")

type Api interface {
	Set(carbonIntensityProvider CarbonIntensityProvider, ctx context.Context) error
	Get(id string, ctx context.Context) (CarbonIntensityProvider, error)
}
