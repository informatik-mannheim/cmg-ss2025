package repo

import (
	"context"
	"encoding/json"
	"os"
	"sync"

	"github.com/informatik-mannheim/cmg-ss2025/services/carbon-intensity-provider/ports"
)

var (
	storageFile     = "zones.json"
	metadataStorage = "zones_metadata.json"
)

type Repo struct {
	carbonIntensityProviders map[string]ports.CarbonIntensityData
	availableZones           []ports.Zone
	mu                       sync.RWMutex
}

// Ensure Repo implements ports.Repo
var _ ports.Repo = (*Repo)(nil)

func NewRepo() *Repo {
	r := &Repo{
		carbonIntensityProviders: make(map[string]ports.CarbonIntensityData),
	}
	r.loadFromFile()
	r.loadZoneMetadata()
	return r
}

// Store a carbon intensity data entry
func (r *Repo) Store(data ports.CarbonIntensityData, ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.carbonIntensityProviders[data.Zone] = data
	return r.saveToFile()
}

// Lookup a specific carbon intensity data entry by zone ID
func (r *Repo) FindById(id string, ctx context.Context) (ports.CarbonIntensityData, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	data, ok := r.carbonIntensityProviders[id]
	if !ok {
		return ports.CarbonIntensityData{}, ports.ErrCarbonIntensityProviderNotFound
	}
	return data, nil
}

// Retrieve all stored carbon intensity data
func (r *Repo) FindAll(ctx context.Context) ([]ports.CarbonIntensityData, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]ports.CarbonIntensityData, 0, len(r.carbonIntensityProviders))
	for _, data := range r.carbonIntensityProviders {
		result = append(result, data)
	}
	return result, nil
}

// Store zone metadata to zones_metadata.json
func (r *Repo) StoreZones(zones []ports.Zone, ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.availableZones = zones
	return r.saveZoneMetadata()
}

// Return zone metadata from memory
func (r *Repo) GetZones(ctx context.Context) []ports.Zone {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.availableZones
}

// Save carbon intensity data to zones.json
func (r *Repo) saveToFile() error {
	file, err := os.Create(storageFile)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(r.carbonIntensityProviders)
}

// Load carbon intensity data from zones.json
func (r *Repo) loadFromFile() {
	file, err := os.Open(storageFile)
	if err != nil {
		return // file doesn't exist, skip
	}
	defer file.Close()
	_ = json.NewDecoder(file).Decode(&r.carbonIntensityProviders)
}

// Save zone metadata to zones_metadata.json
func (r *Repo) saveZoneMetadata() error {
	file, err := os.Create(metadataStorage)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(r.availableZones)
}

// Load zone metadata from zones_metadata.json
func (r *Repo) loadZoneMetadata() {
	file, err := os.Open(metadataStorage)
	if err != nil {
		return
	}
	defer file.Close()
	_ = json.NewDecoder(file).Decode(&r.availableZones)
}
