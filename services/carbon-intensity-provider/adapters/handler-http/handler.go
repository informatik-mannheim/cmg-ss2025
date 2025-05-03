package handler_http

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/informatik-mannheim/cmg-ss2025/services/carbon-intensity-provider/ports"
)

type Handler struct {
	service ports.CarbonIntensityProvider
}

func NewHandler(service ports.CarbonIntensityProvider) http.Handler {
	r := mux.NewRouter()
	h := &Handler{service: service}

	r.HandleFunc("/carbon-intensity/{zone}", h.handleGetCarbonIntensityByZone).Methods("GET")
	r.HandleFunc("/carbon-intensity/zones", h.handleGetAvailableZones).Methods("GET")

	return r
}

func (h *Handler) handleGetCarbonIntensityByZone(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	zone := vars["zone"]

	data, err := h.service.GetCarbonIntensityByZone(zone)
	if err != nil {
		http.Error(w, "Zone not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (h *Handler) handleGetAvailableZones(w http.ResponseWriter, r *http.Request) {
	zones := h.service.GetAvailableZones()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"zones": zones,
	})
}
