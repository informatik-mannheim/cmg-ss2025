package repo

import (
	"context"
	"encoding/json"
	"os"
	"strconv"
	"sync"

	"github.com/informatik-mannheim/cmg-ss2025/services/carbon-intensity-provider/adapters/notifier"
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
	notifier                 notifier.Notifier
}

var _ ports.Repo = (*Repo)(nil)

func NewRepo(n notifier.Notifier) *Repo {
	n.Event("Initializing Repo and loading data from files")

	r := &Repo{
		carbonIntensityProviders: make(map[string]ports.CarbonIntensityData),
		notifier:                 n,
	}

	r.loadFromFile()
	r.loadZoneMetadata()

	return r
}

func (r *Repo) Store(data ports.CarbonIntensityData, ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.carbonIntensityProviders[data.Zone] = data

	if err := r.saveToFile(); err != nil {
		r.notifier.Event("Failed to save carbon intensity data to file: " + err.Error())
		return err
	}

	r.notifier.Event("Stored carbon intensity data for zone: " + data.Zone)
	return nil
}

func (r *Repo) FindById(id string, ctx context.Context) (ports.CarbonIntensityData, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	data, ok := r.carbonIntensityProviders[id]
	if !ok {
		return ports.CarbonIntensityData{}, ports.ErrCarbonIntensityProviderNotFound
	}
	return data, nil
}

func (r *Repo) FindAll(ctx context.Context) ([]ports.CarbonIntensityData, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]ports.CarbonIntensityData, 0, len(r.carbonIntensityProviders))
	for _, data := range r.carbonIntensityProviders {
		result = append(result, data)
	}
	return result, nil
}

func (r *Repo) StoreZones(zones []ports.Zone, ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.availableZones = zones

	if err := r.saveZoneMetadata(); err != nil {
		r.notifier.Event("Failed to save zone metadata: " + err.Error())
		return err
	}

	r.notifier.Event("Stored zone metadata with " + strconv.Itoa(len(zones)) + " zones")
	return nil
}

func (r *Repo) GetZones(ctx context.Context) []ports.Zone {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.availableZones
}

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

func (r *Repo) loadFromFile() {
	file, err := os.Open(storageFile)
	if err != nil {
		r.notifier.Event("No existing carbon intensity data file found (skip load)")
		return
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&r.carbonIntensityProviders); err != nil {
		r.notifier.Event("Failed to decode carbon intensity data: " + err.Error())
	} else {
		r.notifier.Event("Successfully loaded carbon intensity data from file")
	}
}

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

func (r *Repo) loadZoneMetadata() {
	file, err := os.Open(metadataStorage)
	if err != nil {
		r.notifier.Event("No existing zone metadata file found (skip load)")
		return
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&r.availableZones); err != nil {
		r.notifier.Event("Failed to decode zone metadata: " + err.Error())
	} else {
		r.notifier.Event("Successfully loaded zone metadata from file")
	}
}
