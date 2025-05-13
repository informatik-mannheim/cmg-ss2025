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
	Auth       *core.AuthService
	UseLive    bool
	IsAdminFn  func(string) bool
	NotifierFn func() ports.Notifier
}

// New creates a new HTTPHandler
func New(auth *core.AuthService, useLive bool, isAdminFn func(string) bool, notifierFn func() ports.Notifier) *HTTPHandler {
	return &HTTPHandler{
		Auth:       auth,
		UseLive:    useLive,
		IsAdminFn:  isAdminFn,
		NotifierFn: notifierFn,
	}
}

// RegisterHandler handles user registration requests.
func (h *HTTPHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	notifier := h.NotifierFn()
	isAdmin := h.IsAdminFn(r.Header.Get("X-Admin-Secret"))

	var req ports.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil ||
		(req.Role != ports.Consumer && req.Role != ports.Provider && req.Role != ports.JobScheduler) {
		notifier.Event("Invalid register request payload", r.Context())
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.Role == ports.JobScheduler && !isAdmin {
		notifier.Event("Unauthorized attempt to create Job Scheduler", r.Context())
		http.Error(w, "unauthorized to create Job Scheduler", http.StatusForbidden)
		return
	}

	if !isAdmin {
		notifier.Event("Unauthorized register attempt", r.Context())
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id := uuid.NewString()
	notifier.UserRegistered(id, string(req.Role), r.Context())
	notifier.Event("Registration successful for role: "+string(req.Role), r.Context())

	resp := ports.RegisterResponse{ID: id}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)
}

// LoginHandler handles user login requests.
func (h *HTTPHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	notifier := h.NotifierFn()

	var req ports.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Secret == "" {
		notifier.Event("Invalid login request format", r.Context())
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	clientID, token, err := h.Auth.Authenticate(r.Context(), req.Secret)
	if err != nil {
		notifier.Event("Login failed for client: "+clientID, r.Context())
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	notifier.UserLoggedIn(clientID, r.Context())
	notifier.Event("Login successful for client: "+clientID, r.Context())

	resp := ports.LoginResponse{Token: token}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
