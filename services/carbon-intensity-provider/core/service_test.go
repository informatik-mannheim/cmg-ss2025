package core

import (
	"testing"
)

func TestAddAndGetCarbonIntensityByZone(t *testing.T) {
	s := NewCarbonIntensityService()
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
	s := NewCarbonIntensityService()

	_, err := s.GetCarbonIntensityByZone("NOPE")
	if err == nil {
		t.Error("expected error for missing zone, got nil")
	}
}

func TestGetAvailableZones(t *testing.T) {
	s := NewCarbonIntensityService()
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
