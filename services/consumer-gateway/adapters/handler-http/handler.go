package handler_http

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/informatik-mannheim/cmg-ss2025/services/consumer-gateway/ports"
)

type Handler struct {
	job   ports.JobClient
	login ports.LoginClient
	zone  ports.ZoneClient
	rtr   mux.Router
}

var _ http.Handler = (*Handler)(nil)

func NewHandler(job ports.JobClient, login ports.LoginClient, zone ports.ZoneClient) *Handler {
	return &Handler{job: job, login: login, zone: zone}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.rtr.ServeHTTP(w, r) //delegate
}

/*
Creates a new job using the provided data by the client.
The parameter req: contains the fields (imageID, zone) defined
in the CreateJobRequest struct
*/
func (h *Handler) HandleCreateJobRequest(w http.ResponseWriter, r *http.Request) {
	var req ports.CreateJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
		return
	}

	resp, err := h.job.CreateJob(r.Context(), req)
	if err != nil {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

/*
Returns a job result that was requested by client.
The parameter vars: is a map that extracts the path parameters from client request.
So jobs/<job-id>/result returns -> jobID: <job-id>
*/
func (h *Handler) HandleGetJobOutcomeRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobID := vars["job-id"] // "jobID" : "123-abc"

	status, err := h.job.GetJobOutcome(r.Context(), jobID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func (h *Handler) HandleLoginRequest(w http.ResponseWriter, r *http.Request) {
	var req ports.ConsumerLoginRequest // Example: req.Username == "Bob", req.Password == "SuperSecure"
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request"}`, http.StatusUnauthorized)
		return
	}

	resp, err := h.login.Login(r.Context(), req)
	if err != nil {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp) // Token: 123-abc
}
