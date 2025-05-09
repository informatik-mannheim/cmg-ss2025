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

	// Create Job, get available zones and get job result
	h.rtr.HandleFunc("/jobs", h.handleCreateJobRequest).Methods("POST")
	h.rtr.HandleFunc("/carbon-intensity/zones", h.handleGetZones).Methods("GET")
	h.rtr.HandleFunc("/jobs/{id}/result", h.handleGetJobResultRequest).Methods("GET")
	
	// Authentication
	h.rtr.HandleFunc("/auth/login", h.handleLoginRequest).Methods("POST")
	h.rtr.HandleFunc("/auth/register", h.handleRegisterRequest).Methods("POST")
	return &h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.rtr.ServeHTTP(w, r) //delegate
}

/* 
Creates a new job using the provided data by the client.
The parameter req: contains the fields (imageID, zone) defined 
in the CreateJobRequest struct
*/
func (h *Handler) handleCreateJobRequest(w http.ResponseWriter, r *http.Request) {
	var req ports.CreateJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
		return
	}
		
		resp, err := h.service.CreateJob(req, r.Context())
	if err != nil {
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

/*
Returns the available zones from the Carbon Intensity Provider Service.
Choosing a specific location is optional.
*/
func (h *Handler) handleGetZones(w http.ResponseWriter, r *http.Request) {
	var req ports.GetZones
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
		return
	}
		
		resp, err := h.service.GetZones(req, r.Context())
	if err != nil {
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

/* 
Returns a job result that was requested by client.
The parameter vars: is a map that extracts the pathparameters from client request.
So jobs/<job-id>/result returns -> jobID: <job-id>
*/
func (h *Handler) handleGetJobResultRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // 
	jobID := vars["job-id"] // "jobID" : "123-abc"

	status, err := h.service.GetJobResult(jobID, r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}


func (h *Handler) handleLoginRequest(w http.ResponseWriter, r *http.Request) {
	var req ports.ConsumerLoginRequest // Example: req.Username == "Bob", req.Password == "SuperSecure"
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
		return
	}
		
		resp, err := h.service.Login(req, r.Context())
	if err != nil {
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp) // Token: 123-abc
}

func (h *Handler) handleRegisterRequest(w http.ResponseWriter, r *http.Request) {
	var req ports.ConsumerRegistrationRequest 
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
		return
	}
		
		resp, err := h.service.Register(req, r.Context())
	if err != nil {
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
