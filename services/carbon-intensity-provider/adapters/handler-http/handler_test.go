package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	handler "github.com/informatik-mannheim/cmg-ss2025/services/carbon-intensity-provider/adapters/handler-http"
	"github.com/informatik-mannheim/cmg-ss2025/services/carbon-intensity-provider/ports"
)

type mockService struct {
	data  map[string]ports.CarbonIntensityData
	zones []ports.Zone
}

func (m *mockService) GetCarbonIntensityByZone(zone string, ctx context.Context) (ports.CarbonIntensityData, error) {
	if d, ok := m.data[zone]; ok {
		return d, nil
	}
	return ports.CarbonIntensityData{}, ports.ErrCarbonIntensityProviderNotFound
}

func (m *mockService) GetAvailableZones(ctx context.Context) []ports.Zone {
	return m.zones
}

func (m *mockService) AddOrUpdateZone(zone string, intensity float64, ctx context.Context) error {
	return nil
}

func (m *mockService) GetStoredZones(ctx context.Context) []ports.Zone {
	return m.zones
}

type mockNotifier struct {
	events []string
}

func (n *mockNotifier) Event(msg string) {
	n.events = append(n.events, msg)
}

func (n *mockNotifier) CarbonIntensityProviderChanged(data ports.CarbonIntensityData, ctx context.Context) {
	n.events = append(n.events, "notified: "+data.Zone)
}

func TestGetCarbonIntensityByZone(t *testing.T) {
	service := &mockService{
		data: map[string]ports.CarbonIntensityData{
			"GB": {Zone: "GB", CarbonIntensity: 123.4},
		},
	}
	notifier := &mockNotifier{}
	r := handler.NewHandler(service, notifier)

	req := httptest.NewRequest(http.MethodGet, "/carbon-intensity/GB", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	var result ports.CarbonIntensityData
	if err := json.NewDecoder(rr.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if result.Zone != "GB" {
		t.Errorf("Expected zone GB, got %s", result.Zone)
	}
}

func TestGetCarbonIntensityByZone_NotFound(t *testing.T) {
	service := &mockService{data: map[string]ports.CarbonIntensityData{}}
	notifier := &mockNotifier{}
	r := handler.NewHandler(service, notifier)

	req := httptest.NewRequest(http.MethodGet, "/carbon-intensity/XX", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "Zone not found") {
		t.Error("Expected error message in body")
	}
}

func TestGetAvailableZones(t *testing.T) {
	zones := []ports.Zone{
		{Code: "GB", Name: "Great Britain"},
		{Code: "DE", Name: "Germany"},
	}
	service := &mockService{zones: zones}
	notifier := &mockNotifier{}
	r := handler.NewHandler(service, notifier)

	req := httptest.NewRequest(http.MethodGet, "/carbon-intensity/zones", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	var result ports.AvailableZonesResponse
	if err := json.NewDecoder(rr.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if len(result.Zones) != 2 {
		t.Errorf("Expected 2 zones, got %d", len(result.Zones))
	}
}
