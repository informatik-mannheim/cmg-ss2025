package handler_http

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/ports"
)

type Handler struct {
	service ports.Api
	rtr     mux.Router
	secret  string
}

var _ http.Handler = (*Handler)(nil)

func NewHandler(service ports.Api, secret string) *Handler {
	h := Handler{service: service, rtr: *mux.NewRouter(), secret: secret}

	h.rtr.HandleFunc("/schedule", h.handleScheduleJob).Methods("POST")
	h.rtr.HandleFunc("/ping", h.handlePing).Methods("GET")

	return &h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.rtr.ServeHTTP(w, r)
}

func (h *Handler) handleScheduleJob(w http.ResponseWriter, r *http.Request) {
	var reqSecret ports.ScheduleRequest

	if err := json.NewDecoder(r.Body).Decode(&reqSecret); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if reqSecret.Secret != h.secret {
		sendError(w, http.StatusUnauthorized, "Unauthorized: invalid secret")
		return
	}

	if err := h.service.ScheduleJob(); err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to schedule job")
		return
	}

	sendSuccess(w, "Job scheduled successfully")
}

func sendError(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	errorResponse := ports.ScheduleResponse{Error: message}
	if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
		// fallback, better than nothing
		http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
	}
}

func sendSuccess(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusOK)
	successResponse := ports.ScheduleResponse{Message: message}
	if err := json.NewEncoder(w).Encode(successResponse); err != nil {
		http.Error(w, "Failed to encode success response", http.StatusInternalServerError)
	}
}

func (h *Handler) handlePing(w http.ResponseWriter, r *http.Request) {
	sendSuccess(w, "pong")
}
