package handler_http

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/informatik-mannheim/cmg-ss2025/services/consumer-gateway/ports"
)

type Handler struct {
	service ports.Api
	rtr     mux.Router
}

var _ http.Handler = (*Handler)(nil)

func NewHandler(service ports.Api) *Handler {
	h := Handler{service: service, rtr: *mux.NewRouter()}

	h.rtr.HandleFunc("/jobs", h.handleCreateJobRequest).Methods("POST")
	h.rtr.HandleFunc("/jobs/{id}/status", h.handleGetJobStatusResponse).Methods("GET")

	h.rtr.HandleFunc("/auth/login", h.handleLoginRequest).Methods("POST")
	h.rtr.HandleFunc("/auth/register", h.handleRegisterRequest).Methods("POST")
	h.rtr.HandleFunc("/me", h.handleMeRequest).Methods("GET")
	return &h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.rtr.ServeHTTP(w, r) //delegate
}


func (h *Handler) handleGetJobStatusResponse(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // Grabs path parameter
	jobID := vars["job-id"]

	status, err := h.service.GetJobStatus(jobID, r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func (h *Handler) handleCreateJobRequest(w http.ResponseWriter, r *http.Request) {
	var req ports.CreateJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
		return
	}
		
		resp, err := h.service.CreateJobRequest(req, r.Context())
	if err != nil {
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) handleLoginRequest(w http.ResponseWriter, r *http.Request) {
	var req ports.ConsumerLoginRequest // req holds client request data defined in api.go, eg `req.Username == "bob"`` ..
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
		return
	}
		
		resp, err := h.service.ConsumerLoginRequest(req, r.Context())
	if err != nil {
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) handleRegisterRequest(w http.ResponseWriter, r *http.Request) {
	var req ports.ConsumerRegistrationRequest 
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
		return
	}
		
		resp, err := h.service.ConsumerRegisterRequest(req, r.Context())
	if err != nil {
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}


func (h *Handler) handleMeRequest(w http.ResponseWriter, r *http.Request) {
	var req ports.MeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
		return
	}
		
		resp, err := h.service.MeRequest(r.Context())
	if err != nil {
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

