package core

import (
	"context"

	"github.com/informatik-mannheim/cmg-ss2025/services/carbon-intensity-provider/model"
	"github.com/informatik-mannheim/cmg-ss2025/services/carbon-intensity-provider/ports"
)

type CarbonIntensityService struct {
	repo     ports.Repo
	notifier ports.Notifier
}

func NewCarbonIntensityService(repo ports.Repo, notifier ports.Notifier) *CarbonIntensityService {
	return &CarbonIntensityService{
		repo:     repo,
		notifier: notifier,
	}
}

func (s *CarbonIntensityService) GetCarbonIntensityByZone(zone string, ctx context.Context) (model.CarbonIntensityData, error) {
	return s.repo.FindById(zone, ctx)
}

func (s *CarbonIntensityService) GetAvailableZones(ctx context.Context) []model.Zone {
	data, _ := s.repo.FindAll(ctx)

	zones := make([]model.Zone, 0, len(data))
	for _, item := range data {
		zones = append(zones, model.Zone{
			Code: item.Zone,
			Name: item.Zone,
		})
	}
	return zones
}

func (s *CarbonIntensityService) AddOrUpdateZone(zone string, intensity float64, ctx context.Context) error {
	provider := model.CarbonIntensityData{
		Zone:            zone,
		CarbonIntensity: intensity,
	}

	if err := s.repo.Store(provider, ctx); err != nil {
		return err
	}

	// Notifier aufrufen – falls gesetzt
	if s.notifier != nil {
		s.notifier.CarbonIntensityProviderChanged(provider, ctx)
	}

	return nil
}
