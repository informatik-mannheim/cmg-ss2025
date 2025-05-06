package handler

import (
	"bytes"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/informatik-mannheim/cmg-ss2025/services/user-management/core"
	"github.com/informatik-mannheim/cmg-ss2025/services/user-management/model"
	"github.com/informatik-mannheim/cmg-ss2025/services/user-management/notifier"
)

type registerRequest struct {
	Role model.Role `json:"role"`
}

type registerResponse struct {
	ID     string `json:"id"`
	Secret string `json:"secret"`
}

type loginRequest struct {
	Secret string `json:"secret"`
}

type loginResponse struct {
	Token string `json:"token"`
}

var service = core.NewService()

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	isAdmin := isAdmin(r.Header.Get("X-Admin-Secret"))

	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || (req.Role != model.Consumer && req.Role != model.Provider && req.Role != model.JobScheduler) {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.Role == model.JobScheduler && !isAdmin {
		http.Error(w, "unauthorized to create Job Scheduler", http.StatusForbidden)
		return
	}

	if !isAdmin {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id := uuid.NewString()
	secret, err := service.AddUser(id, req.Role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	n := notifier.New()
	n.UserRegistered(id, string(req.Role))

	resp := registerResponse{ID: id, Secret: secret}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Secret == "" {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	user, err := service.Authenticate(req.Secret)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	n := notifier.New()
	n.UserLoggedIn(user.ID)

	token, err := requestAuth0Token(n)
	if err != nil {
		http.Error(w, "failed to retrieve token", http.StatusInternalServerError)
		return
	}

	resp := loginResponse{Token: token}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func requestAuth0Token(n notifier.Notifier) (string, error) {
	type auth0Response struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
	}

	clientID := os.Getenv("AUTH0_CLIENT_ID")
	clientSecret := os.Getenv("AUTH0_CLIENT_SECRET")
	audience := os.Getenv("JWT_AUDIENCE")
	url := os.Getenv("AUTH0_TOKEN_URL")

	n.Event("Requesting Auth0 token from: " + url)
	n.Event("client_id: " + clientID)
	n.Event("audience: " + audience)

	data := map[string]string{
		"client_id":     clientID,
		"client_secret": clientSecret,
		"audience":      audience,
		"grant_type":    "client_credentials",
	}
	jsonData, _ := json.Marshal(data)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		n.Event("Error making request to Auth0: " + err.Error())
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		n.Event("Auth0 responded with status " + resp.Status + ": " + string(body))
		return "", errors.New("Auth0 request failed")
	}

	var parsed auth0Response
	if err := json.Unmarshal(body, &parsed); err != nil {
		n.Event("Failed to parse Auth0 response: " + err.Error())
		return "", err
	}

	n.Event("Successfully received token from Auth0")
	return parsed.AccessToken, nil
}

func isAdmin(headerValue string) bool {
	expectedHash := os.Getenv("ADMIN_SECRET_HASH")
	hash := sha256.Sum256([]byte(headerValue))
	actualHash := hex.EncodeToString(hash[:])
	return subtle.ConstantTimeCompare([]byte(actualHash), []byte(expectedHash)) == 1
}
