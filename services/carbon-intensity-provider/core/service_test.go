package core

import (
	"testing"

	repo_in_memory "github.com/informatik-mannheim/cmg-ss2025/services/carbon-intensity-provider/adapters/repo-in-memory"
)

func TestAddAndGetCarbonIntensityByZone(t *testing.T) {
	repo := repo_in_memory.NewRepo()
	s := NewCarbonIntensityService(repo)
	s.AddOrUpdateZone("DE", 150.0)

	data, err := s.GetCarbonIntensityByZone("DE")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if data.CarbonIntensity != 150.0 {
		t.Errorf("expected intensity 150.0, got: %v", data.CarbonIntensity)
	}
}

func TestGetCarbonIntensityByZone_NotFound(t *testing.T) {
	repo := repo_in_memory.NewRepo()
	s := NewCarbonIntensityService(repo)

	_, err := s.GetCarbonIntensityByZone("NOPE")
	if err == nil {
		t.Error("expected error for missing zone, got nil")
	}
}

func TestGetAvailableZones(t *testing.T) {
	repo := repo_in_memory.NewRepo()
	s := NewCarbonIntensityService(repo)
	s.AddOrUpdateZone("DE", 100.0)
	s.AddOrUpdateZone("FR", 90.0)

	zones := s.GetAvailableZones()
	if len(zones) != 2 {
		t.Errorf("expected 2 zones, got %d", len(zones))
	}

	foundDE, foundFR := false, false
	for _, z := range zones {
		if z.Code == "DE" {
			foundDE = true
		}
		if z.Code == "FR" {
			foundFR = true
		}
	}
	if !foundDE || !foundFR {
		t.Error("expected to find both DE and FR zones")
	}
}

func TestGetAvailableZones_Empty(t *testing.T) {
	repo := repo_in_memory.NewRepo()
	s := NewCarbonIntensityService(repo)

	zones := s.GetAvailableZones()
	if len(zones) != 0 {
		t.Errorf("expected 0 zones, got %d", len(zones))
	}
}
