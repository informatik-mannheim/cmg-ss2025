package ports

import "github.com/informatik-mannheim/cmg-ss2025/services/carbon-intensity-provider/model"

// CarbonIntensityProvider defines the service interface for managing carbon intensity data.
type CarbonIntensityProvider interface {
	GetCarbonIntensityByZone(zone string) (model.CarbonIntensityData, error)
	GetAvailableZones() []model.Zone
}
