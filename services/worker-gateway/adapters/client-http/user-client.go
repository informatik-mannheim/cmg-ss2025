package client_http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/informatik-mannheim/cmg-ss2025/pkg/logging"
	"github.com/informatik-mannheim/cmg-ss2025/services/worker-gateway/ports"
)

type UserClient struct {
	BaseURL    string
	httpClient *http.Client
}

func NewUserClient(baseURL string) *UserClient {
	return &UserClient{BaseURL: baseURL, httpClient: &http.Client{}}
}

func (c *UserClient) GetToken(ctx context.Context) (string, error) {
	url := fmt.Sprintf("%s/auth/register", c.BaseURL)
	reqBody := ports.GetTokenRequest{
		Role: "provider",
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		logging.From(ctx).Error("Failed to marshal request body", "error", err)
		return "", err
	}

	logging.From(ctx).Debug("Registering provider with user service", "url", url)

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		logging.From(ctx).Error("Failed to create HTTP request", "error", err)
		return "", err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		logging.From(ctx).Error("HTTP request failed", "error", err)
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logging.From(ctx).Error("Failed to read response body", "error", err)
		return "", err
	}

	if resp.StatusCode != http.StatusCreated {
		logging.From(ctx).Warn("Unexpected response status from user service", "status", resp.StatusCode, "response", string(respBody))
		return "", fmt.Errorf("user registration failed: %s", string(respBody))
	}

	var parsed ports.GetTokenResponse
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		logging.From(ctx).Error("Failed to unmarshal response", "error", err)
		return "", err
	}

	if parsed.Secret == "" {
		return "", fmt.Errorf("missing 'secret' in response")
	}

	logging.From(ctx).Info("Provider successfully registered", "secret", parsed.Secret)
	return parsed.Secret, nil
}
