package core

import (
	"context"

	"github.com/informatik-mannheim/cmg-ss2025/services/carbon-intensity-provider/ports"
)

type CarbonIntensityService struct {
	repo ports.Repo
}

func NewCarbonIntensityService(repo ports.Repo) *CarbonIntensityService {
	return &CarbonIntensityService{
		repo: repo,
	}
}

func (s *CarbonIntensityService) GetCarbonIntensityByZone(zone string, ctx context.Context) (ports.CarbonIntensityData, error) {
	data, err := s.repo.Fi..ndById(zone, ctx)
	if err != nil {
		return ports.CarbonIntensityData{}, err
	}
	return data, nil
}

func (s *CarbonIntensityService) GetAvailableZones(ctx context.Context) []ports.Zone {
	return s.repo.GetZones(ctx)
}

func (s *CarbonIntensityService) AddOrUpdateZone(zone string, intensity float64, ctx context.Context) error {
	provider := ports.CarbonIntensityData{
		Zone:            zone,
		CarbonIntensity: intensity,
	}

	if err := s.repo.Store(provider, ctx); err != nil {
		return err
	}

	return nil
}

func (s *CarbonIntensityService) GetStoredZones(ctx context.Context) []ports.Zone {
	return s.repo.GetZones(ctx)
}
