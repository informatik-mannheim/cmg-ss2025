package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/informatik-mannheim/cmg-ss2025/services/carbon-intensity-provider/ports"
)

// Handler struct connects HTTP routes to the service logic.
type Handler struct {
	Service ports.CarbonIntensityProvider
}

// NewHandler creates and returns a configured router.
func NewHandler(service ports.CarbonIntensityProvider) *mux.Router {
	h := &Handler{Service: service}
	r := mux.NewRouter()

	r.HandleFunc("/carbon-intensity/zones", h.GetAvailableZones).Methods("GET")
	r.HandleFunc("/carbon-intensity/{zone}", h.GetCarbonIntensityByZone).Methods("GET")

	return r
}

// GetCarbonIntensityByZone handles GET /carbon-intensity/{zone}
func (h *Handler) GetCarbonIntensityByZone(w http.ResponseWriter, r *http.Request) {
	zone := mux.Vars(r)["zone"]
	data, err := h.Service.GetCarbonIntensityByZone(zone, r.Context())
	if err != nil {
		http.Error(w, "Zone not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// GetAvailableZones handles GET /carbon-intensity/zones
func (h *Handler) GetAvailableZones(w http.ResponseWriter, r *http.Request) {
	zoneList := h.Service.GetAvailableZones(r.Context())
	response := ports.AvailableZonesResponse{
		Zones: zoneList,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
