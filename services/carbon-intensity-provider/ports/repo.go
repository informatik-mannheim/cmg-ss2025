package ports

import (
	"context"
	"errors"
)

var ErrCarbonIntensityProviderNotFound = errors.New("carbon intensity provider not found")

type Repo interface {
	Store(data CarbonIntensityData, ctx context.Context) error
	FindById(id string, ctx context.Context) (CarbonIntensityData, error)
	FindAll(ctx context.Context) ([]CarbonIntensityData, error)

	StoreZones([]Zone, context.Context) error
	GetZones(context.Context) []Zone
}
