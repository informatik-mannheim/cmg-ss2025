package core

import (
	"errors"

	"github.com/informatik-mannheim/cmg-ss2025/services/carbon-intensity-provider/model"
)

// CarbonIntensityService is the concrete implementation of the CarbonIntensityProvider interface.
type CarbonIntensityService struct {
	zoneMap map[string]model.CarbonIntensityData
}

// NewCarbonIntensityService creates a new CarbonIntensityService with an empty zone map.
func NewCarbonIntensityService() *CarbonIntensityService {
	return &CarbonIntensityService{
		zoneMap: make(map[string]model.CarbonIntensityData),
	}
}

// GetCarbonIntensityByZone retrieves carbon intensity data for a specific zone.
func (s *CarbonIntensityService) GetCarbonIntensityByZone(zone string) (model.CarbonIntensityData, error) {
	data, exists := s.zoneMap[zone]
	if !exists {
		return model.CarbonIntensityData{}, errors.New("zone not found")
	}
	return data, nil
}

// GetAvailableZones returns a list of all available zones.
func (s *CarbonIntensityService) GetAvailableZones() []model.Zone {
	zones := make([]model.Zone, 0, len(s.zoneMap))
	for code, data := range s.zoneMap {
		zones = append(zones, model.Zone{
			Code: code,
			Name: data.Zone,
		})
	}
	return zones
}

// AddOrUpdateZone adds or updates a zone manually (for testing or seeding).
func (s *CarbonIntensityService) AddOrUpdateZone(zone string, intensity float64) {
	s.zoneMap[zone] = model.CarbonIntensityData{
		Zone:            zone,
		CarbonIntensity: intensity,
	}
}
