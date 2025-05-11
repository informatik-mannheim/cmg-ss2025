package electricitymaps_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	electricitymaps "github.com/informatik-mannheim/cmg-ss2025/services/carbon-intensity-provider/adapters/provider/electricity-maps"
	"github.com/informatik-mannheim/cmg-ss2025/services/carbon-intensity-provider/ports"
)

// MockNotifier is a stub notifier
type MockNotifier struct{}

func (m *MockNotifier) Event(msg string) {}
func (m *MockNotifier) CarbonIntensityProviderChanged(data ports.CarbonIntensityData, ctx context.Context) {
}

var (
	mockNotifier = &MockNotifier{}
	mockTokens   = map[string]string{
		"GB": "test-token",
		"DE": "another-token",
	}
	originalFetchURL        = electricitymaps.FetchURL
	originalZoneMetadataURL = electricitymaps.ZoneMetadataURL
)

func teardown() {
	electricitymaps.FetchURL = originalFetchURL
	electricitymaps.ZoneMetadataURL = originalZoneMetadataURL
}

func newFetcherWithServer(handler http.HandlerFunc) (*electricitymaps.Fetcher, *httptest.Server) {
	server := httptest.NewServer(handler)
	f := electricitymaps.NewWithClient(mockNotifier, server.Client())
	f.TokenByZone = mockTokens
	return f, server
}

func overrideFetchURL(server *httptest.Server) {
	electricitymaps.FetchURL = server.URL + "/carbon-intensity/latest?zone=%s"
}

func overrideZoneURL(server *httptest.Server) {
	electricitymaps.ZoneMetadataURL = server.URL
}

func TestFetch_Success(t *testing.T) {
	defer teardown()

	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("auth-token") != "test-token" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]float64{"carbonIntensity": 111.1})
	}
	f, server := newFetcherWithServer(handler)
	defer server.Close()

	overrideFetchURL(server)

	result, err := f.Fetch("GB", context.Background())
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if result.Zone != "GB" || result.CarbonIntensity != 111.1 {
		t.Errorf("unexpected result: %+v", result)
	}
}

func TestFetch_NoToken(t *testing.T) {
	f := electricitymaps.NewWithClient(mockNotifier, http.DefaultClient)
	_, err := f.Fetch("ZZ", context.Background())
	if err == nil || !strings.Contains(err.Error(), "no token configured") {
		t.Error("expected error due to missing token")
	}
}

func TestFetch_HTTPError(t *testing.T) {
	defer teardown()

	handler := func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "fail", http.StatusInternalServerError)
	}
	f, server := newFetcherWithServer(handler)
	defer server.Close()

	overrideFetchURL(server)

	_, err := f.Fetch("GB", context.Background())
	if err == nil || !strings.Contains(err.Error(), "API returned status") {
		t.Error("expected API status error")
	}
}

func TestFetch_BadJSON(t *testing.T) {
	defer teardown()

	handler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `not-json`)
	}
	f, server := newFetcherWithServer(handler)
	defer server.Close()

	overrideFetchURL(server)

	_, err := f.Fetch("GB", context.Background())
	if err == nil || !strings.Contains(err.Error(), "invalid character") {
		t.Error("expected JSON decode error")
	}
}

func TestNewFromEnv_TokenExtraction(t *testing.T) {
	t.Setenv("TOKEN_GB", "abc")
	t.Setenv("TOKEN_FR", "def")

	f := electricitymaps.NewFromEnv(mockNotifier)

	if f.TokenByZone["GB"] != "abc" {
		t.Error("expected GB token to be 'abc'")
	}
	if f.TokenByZone["FR"] != "def" {
		t.Error("expected FR token to be 'def'")
	}
}

func TestAllElectricityMapZones(t *testing.T) {
	defer teardown()

	mockZones := map[string]map[string]string{
		"GB": {"zoneName": "Great Britain"},
		"DE": {"countryName": "Germany"},
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(mockZones)
	}
	f, server := newFetcherWithServer(handler)
	defer server.Close()

	overrideZoneURL(server)

	zones, err := f.AllElectricityMapZones(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(zones) != 2 {
		t.Fatalf("expected 2 zones, got: %d", len(zones))
	}

	zoneMap := make(map[string]string)
	for _, z := range zones {
		zoneMap[z.Code] = z.Name
	}

	if zoneMap["GB"] != "Great Britain" {
		t.Errorf("expected GB to be 'Great Britain', got: %s", zoneMap["GB"])
	}
	if zoneMap["DE"] != "Germany" {
		t.Errorf("expected DE to be 'Germany', got: %s", zoneMap["DE"])
	}
}

func TestAllElectricityMapZones_BadJSON(t *testing.T) {
	defer teardown()

	handler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "bad json")
	}
	f, server := newFetcherWithServer(handler)
	defer server.Close()

	overrideZoneURL(server)

	_, err := f.AllElectricityMapZones(context.Background())
	if err == nil || !strings.Contains(err.Error(), "invalid character") {
		t.Error("expected JSON decoding error")
	}
}

func TestGetConfiguredZones(t *testing.T) {
	f := electricitymaps.NewWithClient(mockNotifier, http.DefaultClient)
	f.TokenByZone["DE"] = "x"
	f.TokenByZone["FR"] = "y"

	zones, err := f.GetConfiguredZones(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(zones) != 2 {
		t.Errorf("expected 2 zones, got %d", len(zones))
	}
	found := map[string]bool{"DE": false, "FR": false}
	for _, z := range zones {
		found[z] = true
	}
	if !found["DE"] || !found["FR"] {
		t.Error("expected zones DE and FR in configured list")
	}
}

func TestAllElectricityMapZones_BuildRequestError(t *testing.T) {
	electricitymaps.ZoneMetadataURL = "::::::" // invalid URL
	f := electricitymaps.NewWithClient(mockNotifier, http.DefaultClient)

	_, err := f.AllElectricityMapZones(context.Background())
	if err == nil {
		t.Error("expected error due to malformed URL")
	}
}

func TestAllElectricityMapZones_HTTPError(t *testing.T) {
	defer teardown()

	handler := func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "fail", http.StatusInternalServerError)
	}
	f, server := newFetcherWithServer(handler)
	defer server.Close()

	overrideZoneURL(server)

	_, err := f.AllElectricityMapZones(context.Background())
	if err == nil || !strings.Contains(err.Error(), "failed to fetch zones") {
		t.Error("expected fetch zones failure")
	}
}
