package handler_http

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/ports"
)

type Handler struct {
	service ports.Api
	rtr     mux.Router
}

var _ http.Handler = (*Handler)(nil)

func NewHandler(service ports.Api) *Handler {
	h := Handler{service: service, rtr: *mux.NewRouter()}

	h.rtr.HandleFunc("/schedule", h.handleScheduleJob).Methods("POST")

	return &h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.rtr.ServeHTTP(w, r)
}

func (h *Handler) handleScheduleJob(w http.ResponseWriter, r *http.Request) {
	if err := h.service.ScheduleJob(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Job scheduled successfully"))
}
