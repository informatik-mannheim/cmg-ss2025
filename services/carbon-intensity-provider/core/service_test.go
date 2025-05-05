package core

import (
	"context"
	"testing"

	notifier "github.com/informatik-mannheim/cmg-ss2025/services/carbon-intensity-provider/adapters/notifier"
	repo_in_memory "github.com/informatik-mannheim/cmg-ss2025/services/carbon-intensity-provider/adapters/repo-in-memory"
)

func newTestService() *CarbonIntensityService {
	repo := repo_in_memory.NewRepo()
	dummyNotifier := notifier.New()
	return NewCarbonIntensityService(repo, dummyNotifier)
}

func TestAddAndGetCarbonIntensityByZone(t *testing.T) {
	ctx := context.Background()
	s := newTestService()

	err := s.AddOrUpdateZone("DE", 150.0, ctx)
	if err != nil {
		t.Fatalf("unexpected error during AddOrUpdateZone: %v", err)
	}

	data, err := s.GetCarbonIntensityByZone("DE", ctx)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if data.CarbonIntensity != 150.0 {
		t.Errorf("expected intensity 150.0, got: %v", data.CarbonIntensity)
	}
}

func TestGetCarbonIntensityByZone_NotFound(t *testing.T) {
	ctx := context.Background()
	s := newTestService()

	_, err := s.GetCarbonIntensityByZone("NOPE", ctx)
	if err == nil {
		t.Error("expected error for missing zone, got nil")
	}
}

func TestGetAvailableZones(t *testing.T) {
	ctx := context.Background()
	s := newTestService()

	s.AddOrUpdateZone("DE", 100.0, ctx)
	s.AddOrUpdateZone("FR", 90.0, ctx)

	zones := s.GetAvailableZones(ctx)
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
	ctx := context.Background()
	s := newTestService()

	zones := s.GetAvailableZones(ctx)
	if len(zones) != 0 {
		t.Errorf("expected 0 zones, got %d", len(zones))
	}
}
