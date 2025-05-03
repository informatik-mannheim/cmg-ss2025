package core

import (
	"context"

	"github.com/informatik-mannheim/cmg-ss2025/services/carbon-intensity-provider/model"
	"github.com/informatik-mannheim/cmg-ss2025/services/carbon-intensity-provider/ports"
)

// CarbonIntensityService is the concrete implementation of the CarbonIntensityProvider interface.
type CarbonIntensityService struct {
	repo ports.Repo
}

// NewCarbonIntensityService creates a new CarbonIntensityService with an empty zone map.
func NewCarbonIntensityService(repo ports.Repo) *CarbonIntensityService {
	return &CarbonIntensityService{repo: repo}
}

// GetCarbonIntensityByZone retrieves carbon intensity data for a specific zone.
func (s *CarbonIntensityService) GetCarbonIntensityByZone(zone string) (model.CarbonIntensityData, error) {
	return s.repo.FindById(zone, context.Background())
}

// GetAvailableZones returns a list of all available zones.
func (s *CarbonIntensityService) GetAvailableZones() []model.Zone {
	data, _ := s.repo.FindAll(context.Background())

	zones := make([]model.Zone, 0, len(data))
	for _, item := range data {
		zones = append(zones, model.Zone{
			Code: item.Zone,
			Name: item.Zone,
		})
	}
	return zones
}

// AddOrUpdateZone adds or updates a zone manually (for testing or seeding).
func (s *CarbonIntensityService) AddOrUpdateZone(zone string, intensity float64) {
	provider := model.CarbonIntensityData{
		Zone:            zone,
		CarbonIntensity: intensity,
	}
	_ = s.repo.Store(provider, context.Background())
}
