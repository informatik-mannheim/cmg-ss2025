package client_http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/informatik-mannheim/cmg-ss2025/services/consumer-gateway/ports"
)

type LoginClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewLoginClient(baseURL string) *LoginClient {
	return &LoginClient{baseURL: baseURL, httpClient: &http.Client{}}
}

func (c *LoginClient) Login(ctx context.Context, req ports.ConsumerLoginRequest) (ports.LoginResponse, error) {
	url := fmt.Sprintf("%s/auth/login", c.baseURL)

	body, err := json.Marshal(req)
	if err != nil {
		return ports.LoginResponse{}, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return ports.LoginResponse{}, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return ports.LoginResponse{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return ports.LoginResponse{}, fmt.Errorf("user-management login error: %s", resp.Status)
	}

	var out ports.LoginResponse
	err = json.NewDecoder(resp.Body).Decode(&out)
	return out, err
}
