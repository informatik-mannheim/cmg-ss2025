package handler_http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/informatik-mannheim/cmg-ss2025/services/consumer-gateway/ports"
)

type Handler struct {
	api ports.Api
	rtr *mux.Router
}

var _ http.Handler = (*Handler)(nil)

func NewHandler(api ports.Api) *Handler {
	r := mux.NewRouter()
	h := &Handler{api: api, rtr: r}

	r.HandleFunc("/jobs", h.HandleCreateJobRequest).Methods("POST")
	r.HandleFunc("/jobs/{job-id}/outcome", h.HandleGetJobOutcomeRequest).Methods("GET")
	r.HandleFunc("/auth/login", h.HandleLoginRequest).Methods("POST")

	return h
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

	ctx := context.WithValue(r.Context(), "Authorization", r.Header.Get("Authorization"))

	resp, err := h.api.CreateJob(ctx, req)
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

	ctx := context.WithValue(r.Context(), "Authorization", r.Header.Get("Authorization"))

	status, err := h.api.GetJobOutcome(ctx, jobID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func (h *Handler) HandleLoginRequest(w http.ResponseWriter, r *http.Request) {
	var req ports.ConsumerLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request"}`, http.StatusUnauthorized)
		return
	}

	resp, err := h.api.Login(r.Context(), req)
	if err != nil {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
