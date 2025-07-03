package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/informatik-mannheim/cmg-ss2025/services/user-management/core"
	"github.com/informatik-mannheim/cmg-ss2025/services/user-management/ports"
)

// HTTPHandler is an HTTP handler for user management
type HTTPHandler struct {
	Auth      *core.AuthService
	UseLive   bool
	IsAdminFn func(string) bool
}

// New creates a new HTTPHandler
func New(auth *core.AuthService, useLive bool, isAdminFn func(string) bool) *HTTPHandler {
	return &HTTPHandler{
		Auth:      auth,
		UseLive:   useLive,
		IsAdminFn: isAdminFn,
	}
}

// RegisterHandler handles user registration requests.
func (h *HTTPHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	isAdmin := h.IsAdminFn(r.Header.Get("X-Admin-Secret"))

	var req ports.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil ||
		(req.Role != ports.Consumer && req.Role != ports.Provider && req.Role != ports.JobScheduler) {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.Role == ports.JobScheduler && !isAdmin {
		http.Error(w, "unauthorized to create Job Scheduler", http.StatusForbidden)
		return
	}

	if !isAdmin {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	clientID := uuid.NewString()
	clientSecret := uuid.NewString()
	combinedSecret := clientID + "." + clientSecret
	resp := ports.RegisterResponse{Secret: combinedSecret}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)
}

// LoginHandler handles user login requests.
func (h *HTTPHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req ports.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Secret == "" {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	_, token, err := h.Auth.Authenticate(r.Context(), req.Secret)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	resp := ports.LoginResponse{Token: token}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
