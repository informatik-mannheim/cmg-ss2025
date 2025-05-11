package ports

import (
	"context"
)

// CarbonIntensityProvider defines the service interface for managing carbon intensity data.
type CarbonIntensityProvider interface {
	GetCarbonIntensityByZone(zone string, ctx context.Context) (CarbonIntensityData, error)
	GetAvailableZones(ctx context.Context) []Zone
}
