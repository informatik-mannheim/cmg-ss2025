package handler_http

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/informatik-mannheim/cmg-ss2025/services/job/ports"
)

// Handler struct encapsulates the job service and the router
type Handler struct {
	service ports.JobService
	rtr     *mux.Router
}

// Ensure that Handler implements the http.Handler interface
var _ http.Handler = (*Handler)(nil)

// NewHandler initializes the HTTP handlers for each API endpoint
func NewHandler(service ports.JobService) *Handler {
	h := &Handler{service: service, rtr: mux.NewRouter()}
	h.rtr.HandleFunc("/jobs", h.GetJobs).Methods("GET")
	h.rtr.HandleFunc("/jobs", h.CreateJob).Methods("POST")
	h.rtr.HandleFunc("/jobs/{id}", h.GetJob).Methods("GET")
	h.rtr.HandleFunc("/jobs/{id}/outcome", h.GetJobOutcome).Methods("GET")
	h.rtr.HandleFunc("/jobs/{id}/update-scheduler", h.UpdateJobScheduler).Methods("PATCH")
	h.rtr.HandleFunc("/jobs/{id}/update-workerdaemon", h.UpdateJobWorkerDaemon).Methods("PATCH")
	return h
}

// ServeHTTP delegates requests to the appropriate handler in the router
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.rtr.ServeHTTP(w, r)
}

// getJobs handles GET requests to retrieve jobs, possibly filtered by status
func (h *Handler) GetJobs(w http.ResponseWriter, r *http.Request) {
	statusStrings := r.URL.Query()["status"]
	var statuses []ports.JobStatus

	for _, s := range statusStrings {
		parts := strings.Split(s, ",") // Split comma-separated values
		for _, part := range parts {
			part = strings.TrimSpace(part)
			switch ports.JobStatus(part) {
			case ports.StatusQueued, ports.StatusScheduled, ports.StatusRunning,
				ports.StatusCompleted, ports.StatusFailed, ports.StatusCancelled:
				statuses = append(statuses, ports.JobStatus(part))
			default:
				http.Error(w, `{"error": "Bad Request", 
								"message": "Invalid status value: `+part+`"}`, http.StatusBadRequest)
				return
			}
		}
	}

	jobs, err := h.service.GetJobs(r.Context(), statuses)
	if err != nil {
		http.Error(w, HTTPErr500, http.StatusInternalServerError)
		return
	}

	if len(jobs) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jobs)
}

// createJob handles POST requests to create a new job
func (h *Handler) CreateJob(w http.ResponseWriter, r *http.Request) {
	var job ports.JobCreate
	if err := json.NewDecoder(r.Body).Decode(&job); err != nil {
		http.Error(w, HTTPErr400InvalidInputData, http.StatusBadRequest)
		return
	}

	createdJob, err := h.service.CreateJob(r.Context(), job)
	if CheckAndSetErr(w, err) {
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdJob)
}

// getJob retrieves a specific job by ID from the GET request
func (h *Handler) GetJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Check if the ID is a valid UUID
	job, err := h.service.GetJob(r.Context(), id)
	if CheckAndSetErr(w, err) {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(job)
}

// getJobOutcome retrieves the outcome of a job by its ID
func (h *Handler) GetJobOutcome(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	jobOutcome, err := h.service.GetJobOutcome(r.Context(), id)
	if CheckAndSetErr(w, err) {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jobOutcome)
}

// updateJobScheduler handles PATCH requests to update job properties from a scheduler's perspective
func (h *Handler) UpdateJobScheduler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var updateData ports.SchedulerUpdateData
	err := json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		http.Error(w, HTTPErr400InvalidInputData, http.StatusBadRequest)
		return
	}

	updatedJob, err := h.service.UpdateJobScheduler(r.Context(), id, updateData)
	if CheckAndSetErr(w, err) {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedJob)
}

// updateJobWorkerDaemon handles PATCH requests to update job properties from a worker daemon's perspective
func (h *Handler) UpdateJobWorkerDaemon(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var updateData ports.WorkerDaemonUpdateData
	err := json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		http.Error(w, HTTPErr400InvalidInputData, http.StatusBadRequest)
		return
	}

	updatedJob, err := h.service.UpdateJobWorkerDaemon(r.Context(), id, updateData)
	if CheckAndSetErr(w, err) {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedJob)
}

// checkAndSetErr checks for errors and sets the appropriate HTTP response status and message
func CheckAndSetErr(w http.ResponseWriter, err error) bool {
	if err != nil {
		switch err {
		case ports.ErrNotExistingID:
			http.Error(w, HTTPErr400MissId, http.StatusBadRequest)
		case ports.ErrInvalidIDFormat:
			http.Error(w, HTTPErr400InvalidId, http.StatusBadRequest)
		case ports.ErrJobNotFound:
			http.Error(w, HTTPErr400JobNotFound, http.StatusNotFound)
		case ports.ErrNotExistingJobName, ports.ErrNotExistingImageName:
			http.Error(w, HTTPErr400FieldEmpty, http.StatusBadRequest)
		case ports.ErrImageVersionIsInvalid, ports.ErrParamKeyValueEmpty:
			http.Error(w, HTTPErr400InvalidInputData, http.StatusBadRequest)
		case ports.ErrNotExistingStatus:
			http.Error(w, HTTPErr400StatusEmpty, http.StatusBadRequest)
		default:
			http.Error(w, HTTPErr500, http.StatusInternalServerError)
		}
		return true
	}
	return false
}
