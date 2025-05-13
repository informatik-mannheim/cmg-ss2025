package core

import (
	"context"

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

func (s *CarbonIntensityService) GetCarbonIntensityByZone(zone string, ctx context.Context) (ports.CarbonIntensityData, error) {
	data, err := s.repo.FindById(zone, ctx)
	if err != nil {
		s.notifier.Event("Zone not found: " + zone)
		return ports.CarbonIntensityData{}, err
	}

	s.notifier.Event("Retrieved carbon intensity for zone: " + zone)
	s.notifier.CarbonIntensityProviderChanged(data, ctx)
	return data, nil
}

func (s *CarbonIntensityService) GetAvailableZones(ctx context.Context) []ports.Zone {
	zones := s.repo.GetZones(ctx)
	s.notifier.Event("Retrieved all available zones from zone metadata")
	return zones
}

func (s *CarbonIntensityService) AddOrUpdateZone(zone string, intensity float64, ctx context.Context) error {
	provider := ports.CarbonIntensityData{
		Zone:            zone,
		CarbonIntensity: intensity,
	}

	s.notifier.Event("Storing or updating zone: " + zone)

	if err := s.repo.Store(provider, ctx); err != nil {
		s.notifier.Event("Failed to store zone: " + zone)
		return err
	}

	s.notifier.CarbonIntensityProviderChanged(provider, ctx)
	s.notifier.Event("Successfully stored/updated zone: " + zone)

	return nil
}

func (s *CarbonIntensityService) GetStoredZones(ctx context.Context) []ports.Zone {
	return s.repo.GetZones(ctx)
}
