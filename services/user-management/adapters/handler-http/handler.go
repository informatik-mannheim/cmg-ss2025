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

type HTTPHandler struct {
	Auth       ports.AuthProvider
	UseLive    bool
	IsAdminFn  func(string) bool
	NotifierFn func() ports.Notifier
}

type registerRequest struct {
	Role ports.Role `json:"role"`
}

type registerResponse struct {
	ID string `json:"id"`
}

type loginRequest struct {
	Secret string `json:"secret"`
}

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

func (h *HTTPHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	isAdmin := h.IsAdminFn(r.Header.Get("X-Admin-Secret"))

	var req registerRequest
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

	id := uuid.NewString()

	h.NotifierFn().UserRegistered(id, string(req.Role))

	resp := registerResponse{ID: id}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *HTTPHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Secret == "" {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	clientID, _ := splitCredentials(req.Secret)

	token, err := h.Auth.RequestTokenFromCredentials(req.Secret)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	h.NotifierFn().UserLoggedIn(clientID)

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
