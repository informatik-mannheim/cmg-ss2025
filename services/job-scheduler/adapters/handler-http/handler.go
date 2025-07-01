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
	var reqSecret string
	if err := json.NewDecoder(r.Body).Decode(&reqSecret); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if reqSecret != h.secret {
		http.Error(w, "Unauthorized: invalid secret", http.StatusUnauthorized)
		return
	}

	if err := h.service.ScheduleJob(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Job scheduled successfully"))
}

func (h *Handler) handlePing(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}
