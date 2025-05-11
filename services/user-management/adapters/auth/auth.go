package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/informatik-mannheim/cmg-ss2025/services/user-management/ports"
)

type Auth0Adapter struct {
	UseLive  bool
	Notifier ports.Notifier
}

func New(useLive bool, notifier ports.Notifier) *Auth0Adapter {
	return &Auth0Adapter{
		UseLive:  useLive,
		Notifier: notifier,
	}
}

func (a *Auth0Adapter) RequestTokenFromCredentials(ctx context.Context, credentials string) (string, error) {
	a.Notifier.Event("Processing token request from credentials", ctx)

	parts := strings.SplitN(credentials, ".", 2)
	if len(parts) != 2 {
		a.Notifier.Event("Invalid credentials format", ctx)
		return "", fmt.Errorf("invalid credentials format")
	}
	clientID := parts[0]
	clientSecret := parts[1]
	a.Notifier.Event(fmt.Sprintf("Client ID received: %s", clientID), ctx)

	if !a.UseLive {
		a.Notifier.Event("Returning mock token", ctx)
		return mockToken(), nil
	}

	return a.requestRealAuth0Token(ctx, clientID, clientSecret)
}

func (a *Auth0Adapter) requestRealAuth0Token(ctx context.Context, clientID, clientSecret string) (string, error) {
	a.Notifier.Event("Requesting real token from Auth0", ctx)

	type auth0Response struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
	}

	url := os.Getenv("AUTH0_TOKEN_URL")
	audience := os.Getenv("JWT_AUDIENCE")

	a.Notifier.Event(fmt.Sprintf("Using AUTH0_TOKEN_URL: %s", url), ctx)
	a.Notifier.Event(fmt.Sprintf("Using JWT_AUDIENCE: %s", audience), ctx)

	data := map[string]string{
		"client_id":     clientID,
		"client_secret": clientSecret,
		"audience":      audience,
		"grant_type":    "client_credentials",
	}

	body, _ := json.Marshal(data)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		a.Notifier.Event("Failed to create request: "+err.Error(), ctx)
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		a.Notifier.Event("HTTP POST to Auth0 failed: "+err.Error(), ctx)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		msg, _ := io.ReadAll(resp.Body)
		a.Notifier.Event(fmt.Sprintf("Auth0 returned non-200: %d — %s", resp.StatusCode, string(msg)), ctx)
		return "", fmt.Errorf("auth0 failed: %s", msg)
	}

	var parsed auth0Response
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		a.Notifier.Event("Failed to decode Auth0 response: "+err.Error(), ctx)
		return "", err
	}

	a.Notifier.Event("Successfully received token from Auth0", ctx)
	return parsed.AccessToken, nil
}

func mockToken() string {
	return "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IlJGb05QWHZKRmk0SE1vU1lRYzBWNSJ9.eyJodHRwczovL2dyZWVuLWxvYWQtc2hpZnRpbmctcGxhdGZvcm0vcm9sZSI6ImR1bW15X3JvbGUiLCJodHRwczovL2dyZWVuLWxvYWQtc2hpZnRpbmctcGxhdGZvcm0vY2xpZW50X2lkIjoiUWdYSnJrU3Y1WjVkRjhoYzh3cmZPRHYyVk9IZVdCajkiLCJpc3MiOiJodHRwczovL2Rldi1qcWh3Y3U3eHV3Z2RxaTU2LmV1LmF1dGgwLmNvbS8iLCJzdWIiOiJRZ1hKcmtTdjVaNWRGOGhjOHdyZk9EdjJWT0hlV0JqOUBjbGllbnRzIiwiYXVkIjoiaHR0cHM6Ly9ncmVlbi1sb2FkLXNoaWZ0aW5nLXBsYXRmb3JtLyIsImlhdCI6MTc0NjkxMTExNiwiZXhwIjoxNzQ2OTk3NTE2LCJndHkiOiJjbGllbnQtY3JlZGVudGlhbHMiLCJhenAiOiJRZ1hKcmtTdjVaNWRGOGhjOHdyZk9EdjJWT0hlV0JqOSIsInBlcm1pc3Npb25zIjpbXX0.sVRxi8Ea-_3GhkTUOhtDH8Io8Ds3u-TiYELq2wtGwVE4iFzFXRdRwEKSIEz6ELCt_MVjYVlvdza1hnQdgSKmxOp_Hs7ZKnCYRqrFEyKff0_kzHKWR65e0gpBniMKhh97vZ8jmTWOf7F39nIJCZNZ3RFyrkiXCvyxQKXujmJnfRlXKbr9AdRVQGFL-QtDEVSstG_b0954J1zhCAp3dUOSbvo3h1TQI0sZz_WNQOOSaWaH0m9oTzdMdOXkvOfqD3A7Zw8cihwxzITjLKAjVA276wZmcFXe-E5o45uXV5nPaDl4GxPZa3fwFNUB6h9-EL3uPFXRx9RZdr2hmthP-4vcJQ"
}
