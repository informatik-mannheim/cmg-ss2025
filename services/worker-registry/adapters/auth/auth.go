package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/informatik-mannheim/cmg-ss2025/services/worker-registry/ports"
)

func GetAuthEndpoint(base string) string {
	return fmt.Sprintf("%s/auth/login", base)
}

type AuthAdapter struct {
	baseUrl string
	secret  string
	token   string
}

var _ ports.AuthAdapter = (*AuthAdapter)(nil)

func NewAuthAdapter(baseUrl, secret string) *AuthAdapter {
	return &AuthAdapter{
		baseUrl: baseUrl,
		secret:  secret,
	}
}

func (adapter *AuthAdapter) Authenticate() error {
	endpoint := GetAuthEndpoint(adapter.baseUrl)

	data := ports.GetAuthToken{
		Secret: adapter.secret,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal auth request: %w", err)
	}
	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send POST request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("authentication failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var result ports.AuthTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	adapter.token = result.Token
	return nil
}

func (adapter *AuthAdapter) GetToken() string {
	return adapter.token
}
