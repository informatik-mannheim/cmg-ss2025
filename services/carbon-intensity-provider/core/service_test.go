package core_test

import (
	"context"
	"testing"

	"github.com/informatik-mannheim/cmg-ss2025/services/carbon-intensity-provider/core"
	"github.com/informatik-mannheim/cmg-ss2025/services/carbon-intensity-provider/ports"
)

// MockRepo implements ports.Repo
type MockRepo struct {
	storage  map[string]ports.CarbonIntensityData
	zones    []ports.Zone
	storeErr error
}

func (m *MockRepo) Store(data ports.CarbonIntensityData, ctx context.Context) error {
	if m.storage == nil {
		m.storage = make(map[string]ports.CarbonIntensityData)
	}
	if m.storeErr != nil {
		return m.storeErr
	}
	m.storage[data.Zone] = data
	return nil
}

func (m *MockRepo) FindById(id string, ctx context.Context) (ports.CarbonIntensityData, error) {
	if m.storage == nil {
		return ports.CarbonIntensityData{}, ports.ErrCarbonIntensityProviderNotFound
	}
	data, ok := m.storage[id]
	if !ok {
		return ports.CarbonIntensityData{}, ports.ErrCarbonIntensityProviderNotFound
	}
	return data, nil
}

func (m *MockRepo) FindAll(ctx context.Context) ([]ports.CarbonIntensityData, error) {
	var result []ports.CarbonIntensityData
	for _, v := range m.storage {
		result = append(result, v)
	}
	return result, nil
}

func (m *MockRepo) StoreZones(zones []ports.Zone, ctx context.Context) error {
	m.zones = zones
	return nil
}

func (m *MockRepo) GetZones(ctx context.Context) []ports.Zone {
	return m.zones
}

// MockNotifier implements ports.Notifier
func TestAddOrUpdateZone_Success(t *testing.T) {
	repo := &MockRepo{}
	service := core.NewCarbonIntensityService(repo)

	err := service.AddOrUpdateZone("DE", 100.0, context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if val, ok := repo.storage["DE"]; !ok || val.CarbonIntensity != 100.0 {
		t.Errorf("expected stored value for DE to be 100.0, got %+v", val)
	}
}

func TestGetCarbonIntensityByZone_Found(t *testing.T) {
	repo := &MockRepo{
		storage: map[string]ports.CarbonIntensityData{
			"FR": {Zone: "FR", CarbonIntensity: 90.0},
		},
	}
	service := core.NewCarbonIntensityService(repo)

	data, err := service.GetCarbonIntensityByZone("FR", context.Background())
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if data.CarbonIntensity != 90.0 {
		t.Errorf("expected 90.0, got %.2f", data.CarbonIntensity)
	}
}

func TestGetCarbonIntensityByZone_NotFound(t *testing.T) {
	repo := &MockRepo{}
	service := core.NewCarbonIntensityService(repo)

	_, err := service.GetCarbonIntensityByZone("NOPE", context.Background())
	if err == nil {
		t.Error("expected error for unknown zone")
	}
}

func TestGetAvailableZones(t *testing.T) {
	repo := &MockRepo{
		zones: []ports.Zone{
			{Code: "DE", Name: "Germany"},
			{Code: "FR", Name: "France"},
		},
	}
	service := core.NewCarbonIntensityService(repo)

	zones := service.GetAvailableZones(context.Background())
	if len(zones) != 2 {
		t.Errorf("expected 2 zones, got %d", len(zones))
	}
}
