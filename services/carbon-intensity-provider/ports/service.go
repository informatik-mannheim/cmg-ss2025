package ports

import (
	"context"

	"github.com/informatik-mannheim/cmg-ss2025/services/carbon-intensity-provider/model"
)

// CarbonIntensityProvider defines the service interface for managing carbon intensity data.
type CarbonIntensityProvider interface {
	GetCarbonIntensityByZone(zone string, ctx context.Context) (model.CarbonIntensityData, error)
	GetAvailableZones(ctx context.Context) []model.Zone
}
