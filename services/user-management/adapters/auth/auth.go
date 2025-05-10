package auth0

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type Auth0Adapter struct {
	UseLive bool
}

func New(useLive bool) *Auth0Adapter {
	return &Auth0Adapter{UseLive: useLive}
}

func (a *Auth0Adapter) RequestTokenFromCredentials(credentials string) (string, error) {
	parts := strings.SplitN(credentials, ".", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid credentials format")
	}
	clientID := parts[0]
	clientSecret := parts[1]

	if !a.UseLive {
		return "local-jwt-for-" + clientID, nil
	}

	return a.requestRealAuth0Token(clientID, clientSecret)
}

func (a *Auth0Adapter) requestRealAuth0Token(clientID, clientSecret string) (string, error) {
	type auth0Response struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
	}

	data := map[string]string{
		"client_id":     clientID,
		"client_secret": clientSecret,
		"audience":      os.Getenv("JWT_AUDIENCE"),
		"grant_type":    "client_credentials",
	}

	body, _ := json.Marshal(data)
	resp, err := http.Post(os.Getenv("AUTH0_TOKEN_URL"), "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		msg, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("auth0 failed: %s", msg)
	}

	var parsed auth0Response
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return "", err
	}
	return parsed.AccessToken, nil
}
