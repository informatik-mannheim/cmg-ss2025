package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Auth0Adapter struct {
	UseLive bool
}

// New creates a new Auth0Adapter with optional live mode
func New(useLive bool) *Auth0Adapter {
	return &Auth0Adapter{
		UseLive: useLive,
	}
}

// RequestTokenFromClientSecret implements ports.TokenProvider
func (a *Auth0Adapter) RequestTokenFromClientSecret(ctx context.Context, clientID, clientSecret string) (string, error) {
	if !a.UseLive {
		return mockToken(), nil
	}

	url := os.Getenv("AUTH0_TOKEN_URL")
	audience := os.Getenv("JWT_AUDIENCE")

	data := map[string]string{
		"client_id":     clientID,
		"client_secret": clientSecret,
		"audience":      audience,
		"grant_type":    "client_credentials",
	}

	body, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		msg, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("auth0 error: %s", msg)
	}

	var parsed struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return "", err
	}

	return parsed.AccessToken, nil
}

func mockToken() string {
	return "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IlJGb05QWHZKRmk0SE1vU1lRYzBWNSJ9.eyJodHRwczovL2dyZWVuLWxvYWQtc2hpZnRpbmctcGxhdGZvcm0vcm9sZSI6ImR1bW15X3JvbGUiLCJodHRwczovL2dyZWVuLWxvYWQtc2hpZnRpbmctcGxhdGZvcm0vY2xpZW50X2lkIjoiUWdYSnJrU3Y1WjVkRjhoYzh3cmZPRHYyVk9IZVdCajkiLCJpc3MiOiJodHRwczovL2Rldi1qcWh3Y3U3eHV3Z2RxaTU2LmV1LmF1dGgwLmNvbS8iLCJzdWIiOiJRZ1hKcmtTdjVaNWRGOGhjOHdyZk9EdjJWT0hlV0JqOUBjbGllbnRzIiwiYXVkIjoiaHR0cHM6Ly9ncmVlbi1sb2FkLXNoaWZ0aW5nLXBsYXRmb3JtLyIsImlhdCI6MTc0NjkxMTExNiwiZXhwIjoxNzQ2OTk3NTE2LCJndHkiOiJjbGllbnQtY3JlZGVudGlhbHMiLCJhenAiOiJRZ1hKcmtTdjVaNWRGOGhjOHdyZk9EdjJWT0hlV0JqOSIsInBlcm1pc3Npb25zIjpbXX0.sVRxi8Ea-_3GhkTUOhtDH8Io8Ds3u-TiYELq2wtGwVE4iFzFXRdRwEKSIEz6ELCt_MVjYVlvdza1hnQdgSKmxOp_Hs7ZKnCYRqrFEyKff0_kzHKWR65e0gpBniMKhh97vZ8jmTWOf7F39nIJCZNZ3RFyrkiXCvyxQKXujmJnfRlXKbr9AdRVQGFL-QtDEVSstG_b0954J1zhCAp3dUOSbvo3h1TQI0sZz_WNQOOSaWaH0m9oTzdMdOXkvOfqD3A7Zw8cihwxzITjLKAjVA276wZmcFXe-E5o45uXV5nPaDl4GxPZa3fwFNUB6h9-EL3uPFXRx9RZdr2hmthP-4vcJQ"
}
