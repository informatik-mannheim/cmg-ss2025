package repo

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/informatik-mannheim/cmg-ss2025/services/carbon-intensity-provider/ports"
)

func setupTestRepo(t *testing.T) *Repo {
	t.Helper()

	// Override storage paths to temp dir
	tmpDir := t.TempDir()
	storageFile = filepath.Join(tmpDir, "zones.json")
	metadataStorage = filepath.Join(tmpDir, "zones_metadata.json")

	return NewRepo()
}

func TestRepo_StoreAndFind(t *testing.T) {
	ctx := context.Background()
	r := setupTestRepo(t)

	data := ports.CarbonIntensityData{Zone: "DE", CarbonIntensity: 150.0}
	if err := r.Store(data, ctx); err != nil {
		t.Fatalf("expected no error on Store, got %v", err)
	}

	got, err := r.FindById("DE", ctx)
	if err != nil {
		t.Fatalf("expected to find zone DE, got error: %v", err)
	}
	if got.CarbonIntensity != 150.0 {
		t.Errorf("expected intensity 150.0, got %.2f", got.CarbonIntensity)
	}
}

func TestRepo_FindAll(t *testing.T) {
	ctx := context.Background()
	r := setupTestRepo(t)

	_ = r.Store(ports.CarbonIntensityData{Zone: "DE", CarbonIntensity: 150.0}, ctx)
	_ = r.Store(ports.CarbonIntensityData{Zone: "FR", CarbonIntensity: 90.0}, ctx)

	all, err := r.FindAll(ctx)
	if err != nil {
		t.Fatalf("FindAll error: %v", err)
	}
	if len(all) < 2 {
		t.Errorf("expected at least 2 entries, got %d", len(all))
	}
}

func TestRepo_FindById_NotFound(t *testing.T) {
	ctx := context.Background()
	r := setupTestRepo(t)

	_, err := r.FindById("NOPE", ctx)
	if err == nil {
		t.Error("expected error for unknown ID, got nil")
	}
}

func TestRepo_StoreAndGetZones(t *testing.T) {
	ctx := context.Background()
	r := setupTestRepo(t)

	zones := []ports.Zone{
		{Code: "DE", Name: "Germany"},
		{Code: "FR", Name: "France"},
	}
	if err := r.StoreZones(zones, ctx); err != nil {
		t.Fatalf("failed to store zones: %v", err)
	}

	stored := r.GetZones(ctx)
	if len(stored) != 2 {
		t.Errorf("expected 2 stored zones, got %d", len(stored))
	}
}

func TestRepo_PersistenceFilesCreated(t *testing.T) {
	r := setupTestRepo(t)
	ctx := context.Background()

	_ = r.Store(ports.CarbonIntensityData{Zone: "DE", CarbonIntensity: 100.0}, ctx)
	_ = r.StoreZones([]ports.Zone{{Code: "DE", Name: "Germany"}}, ctx)

	if _, err := os.Stat(storageFile); err != nil {
		t.Errorf("zones.json not created: %v", err)
	}
	if _, err := os.Stat(metadataStorage); err != nil {
		t.Errorf("zones_metadata.json not created: %v", err)
	}
}
