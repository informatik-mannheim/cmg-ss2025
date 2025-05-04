package handler_http

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/informatik-mannheim/cmg-ss2025/services/worker-registry/ports"
)

type Handler struct {
	service ports.Api
	rtr     mux.Router
}

var _ http.Handler = (*Handler)(nil)

func NewHandler(service ports.Api) *Handler {
	h := Handler{service: service, rtr: *mux.NewRouter()}

	h.rtr.HandleFunc("/workers", h.handleGetAll).Methods("GET")
	h.rtr.HandleFunc("/workers", h.handleCreate).Methods("POST")
	h.rtr.HandleFunc("/workers/{id}", h.handleGetById).Methods("GET")
	h.rtr.HandleFunc("/workers/{id}/status", h.handleUpdateStatus).Methods("PUT")

	return &h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.rtr.ServeHTTP(w, r)
}

func (h *Handler) handleGetAll(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	zone := r.URL.Query().Get("zone")

	workers, err := h.service.GetWorkers(status, zone, r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workers)
}

func (h *Handler) handleGetById(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	worker, err := h.service.GetWorkerById(id, r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(worker)
}

func (h *Handler) handleCreate(w http.ResponseWriter, r *http.Request) {
	zone := r.URL.Query().Get("zone")

	if zone == "" {
		http.Error(w, "zone is required to create Workers", http.StatusBadRequest)
		return
	}

	worker, err := h.service.CreateWorker(zone, r.Context())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(worker)
}

func (h *Handler) handleUpdateStatus(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var payload struct {
		Status string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updatedWorker, err := h.service.UpdateWorkerStatus(id, payload.Status, r.Context())

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedWorker)
	w.WriteHeader(http.StatusOK)
}
