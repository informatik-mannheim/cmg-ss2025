package handler_http

import (
	"encoding/json"
	"net/http"

	"github.com/informatik-mannheim/cmg-ss2025/services/worker-gateway/ports"
)

type Handler struct {
	api ports.Api
}

func NewHandler(api ports.Api) *Handler {
	return &Handler{api: api}
}

// POST /worker/heartbeat
func (h *Handler) HeartbeatHandler(w http.ResponseWriter, r *http.Request) {
	var req ports.HeartbeatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	jobs, err := h.api.Heartbeat(r.Context(), req)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jobs)
}

// POST /result
func (h *Handler) SubmitResultHandler(w http.ResponseWriter, r *http.Request) {
	var result ports.ResultRequest
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if err := h.api.Result(r.Context(), result); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// POST /register
func (h *Handler) RegisterWorkerHandler(w http.ResponseWriter, r *http.Request) {
	var req ports.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if err := h.api.Register(r.Context(), req); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
