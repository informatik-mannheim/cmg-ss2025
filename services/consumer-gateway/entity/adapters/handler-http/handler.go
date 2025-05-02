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
	h.rtr.HandleFunc("/jobs/{id}/status", h.handleGetJobResultRequest).Methods("GET")

	h.rtr.HandleFunc("/auth/login", h.handleLoginRequest).Methods("POST")
	h.rtr.HandleFunc("/auth/register", h.handleRegisterRequest).Methods("POST")
	return &h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.rtr.ServeHTTP(w, r) //delegate
}


/* 
Creates a new job using the provided data ba the client.
The parameter req: contains the fields (imageID, location) defined in 
the CreateJobRequest struct
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
Returns a job status that was requested by client.
The parameter vars: is a map that extracts the pathparameters from client request.
So jobs/<job-id>/status returns -> jobID: <job-id>
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
