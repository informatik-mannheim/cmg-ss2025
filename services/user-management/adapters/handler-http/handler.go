// handler/handler.go
package handler

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/informatik-mannheim/cmg-ss2025/services/user-management/ports"
)

// HTTPHandler is an HTTP handler for user management
// It provides endpoints for user registration and login
// It uses the AuthProvider interface to handle authentication
// It uses the Notifier interface to send notifications
// It can be configured to use a live Auth0 instance
// or to use a mock token for testing purposes
type HTTPHandler struct {
	Auth       ports.AuthProvider
	UseLive    bool
	IsAdminFn  func(string) bool
	NotifierFn func() ports.Notifier
}

// registerRequest is the request payload for user registration.
// It contains the role of the user to be registered.
// The role can be either Consumer, Provider, or JobScheduler.
type registerRequest struct {
	Role ports.Role `json:"role"`
}

// registerResponse is the response payload for user registration.
// It contains the generated user ID.
type registerResponse struct {
	ID string `json:"id"`
}

// loginRequest is the request payload for user login.
// It contains the secret used for authentication.
// The secret should be in the format "clientID.clientSecret".
type loginRequest struct {
	Secret string `json:"secret"`
}

// loginResponse is the response payload for user login.
// It contains the token received from Auth0.
type loginResponse struct {
	Token string `json:"token"`
}

func New(auth ports.AuthProvider, useLive bool, isAdminFn func(string) bool, notifierFn func() ports.Notifier) *HTTPHandler {
	return &HTTPHandler{
		Auth:       auth,
		UseLive:    useLive,
		IsAdminFn:  isAdminFn,
		NotifierFn: notifierFn,
	}
}

// RegisterHandler handles user registration requests.
// It checks the request payload for validity and role.
// If the request is valid, it generates a new user ID and sends a notification.
// It returns a 201 Created response with the user ID.
// If the request is invalid, it returns a 400 Bad Request response.
// If the user is not authorized to create a Job Scheduler, it returns a 403 Forbidden response.
// If the user is not authorized to register, it returns a 401 Unauthorized response.
func (h *HTTPHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	notifier := h.NotifierFn()
	isAdmin := h.IsAdminFn(r.Header.Get("X-Admin-Secret"))

	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil ||
		(req.Role != ports.Consumer && req.Role != ports.Provider && req.Role != ports.JobScheduler) {
		notifier.Event("Invalid register request payload")
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.Role == ports.JobScheduler && !isAdmin {
		notifier.Event("Unauthorized attempt to create Job Scheduler")
		http.Error(w, "unauthorized to create Job Scheduler", http.StatusForbidden)
		return
	}

	if !isAdmin {
		notifier.Event("Unauthorized register attempt")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id := uuid.NewString()
	notifier.UserRegistered(id, string(req.Role))
	notifier.Event("Registration successful for role: " + string(req.Role))

	resp := registerResponse{ID: id}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// LoginHandler handles user login requests.
// It checks the request payload for validity and secret.
func (h *HTTPHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	notifier := h.NotifierFn()

	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Secret == "" {
		notifier.Event("Invalid login request format")
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	clientID, _ := splitCredentials(req.Secret)
	notifier.Event("Login attempt from client: " + clientID)
	token, err := h.Auth.RequestTokenFromCredentials(req.Secret)
	if err != nil {
		notifier.Event("Login failed for client: " + clientID)
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	notifier.UserLoggedIn(clientID)
	notifier.Event("Login successful for client: " + clientID)

	resp := loginResponse{Token: token}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func IsAdmin(headerValue string) bool {
	expectedHash := os.Getenv("ADMIN_SECRET_HASH")
	hash := sha256.Sum256([]byte(headerValue))
	actualHash := hex.EncodeToString(hash[:])
	return subtle.ConstantTimeCompare([]byte(actualHash), []byte(expectedHash)) == 1
}

func splitCredentials(secret string) (string, string) {
	parts := strings.SplitN(secret, ".", 2)
	if len(parts) != 2 {
		return "", ""
	}
	return parts[0], parts[1]
}
