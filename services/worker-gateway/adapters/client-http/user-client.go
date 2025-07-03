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

func (c *UserClient) GetToken(ctx context.Context, req ports.GetTokenRequest) (ports.GetTokenResponse, error) {
	url := fmt.Sprintf("%s/auth/login", c.BaseURL)

	bodyBytes, err := json.Marshal(req)
	if err != nil {
		logging.From(ctx).Error("Failed to marshal request body", "error", err)
		return ports.GetTokenResponse{}, err
	}

	logging.From(ctx).Debug("Login Worker with user service", "url", url)

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		logging.From(ctx).Error("Failed to create HTTP request", "error", err)
		return ports.GetTokenResponse{}, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		logging.From(ctx).Error("HTTP request failed", "error", err)
		return ports.GetTokenResponse{}, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logging.From(ctx).Error("Failed to read response body", "error", err)
		return ports.GetTokenResponse{}, err
	}

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		logging.From(ctx).Warn("Unexpected response status from user service", "status", resp.StatusCode, "response", string(respBody))
		return ports.GetTokenResponse{}, fmt.Errorf("user registration failed: %s", string(respBody))
	}

	var parsed ports.GetTokenResponse
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		logging.From(ctx).Error("Failed to unmarshal response", "error", err)
		return ports.GetTokenResponse{}, err
	}

	if parsed.Token == "" {
		return ports.GetTokenResponse{}, fmt.Errorf("missing 'token' in response")
	}

	logging.From(ctx).Debug("Provider successfully registered")
	return parsed, nil
}
