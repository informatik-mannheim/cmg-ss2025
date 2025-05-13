package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/informatik-mannheim/cmg-ss2025/services/carbon-intensity-provider/ports"
)

// Handler struct connects HTTP routes to the service logic.
type Handler struct {
	Service  ports.CarbonIntensityProvider
	Notifier ports.Notifier
}

// NewHandler creates and returns a configured router.
func NewHandler(service ports.CarbonIntensityProvider, notifier ports.Notifier) *mux.Router {
	h := &Handler{
		Service:  service,
		Notifier: notifier,
	}

	r := mux.NewRouter()
	r.HandleFunc("/carbon-intensity/zones", h.GetAvailableZones).Methods("GET")
	r.HandleFunc("/carbon-intensity/{zone}", h.GetCarbonIntensityByZone).Methods("GET")

	return r
}

// GetCarbonIntensityByZone handles GET /carbon-intensity/{zone}
func (h *Handler) GetCarbonIntensityByZone(w http.ResponseWriter, r *http.Request) {
	zone := mux.Vars(r)["zone"]
	h.Notifier.Event("API: GET /carbon-intensity/" + zone)

	data, err := h.Service.GetCarbonIntensityByZone(zone, r.Context())
	if err != nil {
		h.Notifier.Event("Zone not found: " + zone)
		http.Error(w, "Zone not found", http.StatusNotFound)
		return
	}

	h.Notifier.Event("Returning carbon data for zone: " + zone)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// GetAvailableZones handles GET /carbon-intensity/zones
func (h *Handler) GetAvailableZones(w http.ResponseWriter, r *http.Request) {
	h.Notifier.Event("API: GET /carbon-intensity/zones")

	zones := h.Service.GetAvailableZones(r.Context())
	response := ports.AvailableZonesResponse{
		Zones: zones,
	}

	h.Notifier.Event("Returning stored zone list with " + strconv.Itoa(len(response.Zones)) + " zones")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
