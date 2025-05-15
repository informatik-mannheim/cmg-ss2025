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
	h.rtr.HandleFunc("/jobs", h.getJobs).Methods("GET")
	h.rtr.HandleFunc("/jobs", h.createJob).Methods("POST")
	h.rtr.HandleFunc("/jobs/{id}", h.getJob).Methods("GET")
	h.rtr.HandleFunc("/jobs/{id}/outcome", h.getJobOutcome).Methods("GET")
	h.rtr.HandleFunc("/jobs/{id}/update-scheduler", h.updateJobScheduler).Methods("PATCH")
	h.rtr.HandleFunc("/jobs/{id}/update-workerdaemon", h.updateJobWorkerDaemon).Methods("PATCH")
	return h
}

// ServeHTTP delegates requests to the appropriate handler in the router
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.rtr.ServeHTTP(w, r)
}

// getJobs handles GET requests to retrieve jobs, possibly filtered by status
func (h *Handler) getJobs(w http.ResponseWriter, r *http.Request) {
	statusStrings := r.URL.Query()["status"]
	var statuses []ports.JobStatus

	for _, s := range statusStrings {
		parts := strings.Split(s, ",") // Split comma-separated values
		for _, part := range parts {
			part = strings.TrimSpace(part)
			switch part {
			case "queued", "scheduled", "running", "completed", "failed", "cancelled":
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
		http.Error(w, ports.HTTPErr500, http.StatusInternalServerError)
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
func (h *Handler) createJob(w http.ResponseWriter, r *http.Request) {
	var job ports.JobCreate
	err := json.NewDecoder(r.Body).Decode(&job)
	if err != nil {
		http.Error(w, ports.HTTPErr400InvalidInputData, http.StatusBadRequest)
		return
	}

	// Validate required fields in `job` (adjust validation according to your struct and requirements)
	if job.JobName == "" || job.CreationZone == "" || job.Image.Name == "" {
		http.Error(w, ports.HTTPErr400FieldEmpty, http.StatusBadRequest)
		return
	}

	createdJob, err := h.service.CreateJob(r.Context(), job)
	if err != nil {
		http.Error(w, ports.HTTPErr500, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdJob)
}

// getJob retrieves a specific job by ID from the GET request
func (h *Handler) getJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Check if the ID is a valid UUID
	job, err := h.service.GetJob(r.Context(), id)
	if err != nil {
		if err == ports.ErrNotExistingID {
			http.Error(w, ports.HTTPErr400MissId, http.StatusBadRequest)
		} else if err == ports.ErrInvalidIDFormat {
			http.Error(w, ports.HTTPErr400InvalidId, http.StatusBadRequest)
		} else if err == ports.ErrJobNotFound {
			http.Error(w, ports.HTTPErr400JobNotFound, http.StatusNotFound)
		} else {
			http.Error(w, ports.HTTPErr500, http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(job)
}

// getJobOutcome retrieves the outcome of a job by its ID
func (h *Handler) getJobOutcome(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	jobOutcome, err := h.service.GetJobOutcome(r.Context(), id)
	if err != nil {
		if err == ports.ErrNotExistingID {
			http.Error(w, ports.HTTPErr400MissId, http.StatusBadRequest)
		} else if err == ports.ErrInvalidIDFormat {
			http.Error(w, ports.HTTPErr400InvalidId, http.StatusBadRequest)
		} else if err == ports.ErrJobNotFound {
			http.Error(w, ports.HTTPErr400JobNotFound, http.StatusNotFound)
		} else {
			http.Error(w, ports.HTTPErr500, http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jobOutcome)
}

// updateJobScheduler handles PATCH requests to update job properties from a scheduler's perspective
func (h *Handler) updateJobScheduler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var updateData ports.SchedulerUpdateData
	err := json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		http.Error(w, ports.HTTPErr400InvalidInputData, http.StatusBadRequest)
		return
	}

	updatedJob, err := h.service.UpdateJobScheduler(r.Context(), id, updateData)
	if err != nil {
		if err == ports.ErrNotExistingID {
			http.Error(w, ports.HTTPErr400MissId, http.StatusBadRequest)
		} else if err == ports.ErrInvalidIDFormat {
			http.Error(w, ports.HTTPErr400InvalidId, http.StatusBadRequest)
		} else if err == ports.ErrJobNotFound {
			http.Error(w, ports.HTTPErr400JobNotFound, http.StatusNotFound)
		} else {
			http.Error(w, ports.HTTPErr500, http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedJob)
}

// updateJobWorkerDaemon handles PATCH requests to update job properties from a worker daemon's perspective
func (h *Handler) updateJobWorkerDaemon(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var updateData ports.WorkerDaemonUpdateData
	err := json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		http.Error(w, ports.HTTPErr400InvalidInputData, http.StatusBadRequest)
		return
	}

	updatedJob, err := h.service.UpdateJobWorkerDaemon(r.Context(), id, updateData)
	if err != nil {
		if err == ports.ErrNotExistingID {
			http.Error(w, ports.HTTPErr400MissId, http.StatusBadRequest)
		} else if err == ports.ErrInvalidIDFormat {
			http.Error(w, ports.HTTPErr400InvalidId, http.StatusBadRequest)
		} else if err == ports.ErrJobNotFound {
			http.Error(w, ports.HTTPErr400JobNotFound, http.StatusNotFound)
		} else {
			http.Error(w, ports.HTTPErr500, http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedJob)
}
