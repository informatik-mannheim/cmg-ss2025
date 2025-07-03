package core_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"testing"

	"github.com/informatik-mannheim/cmg-ss2025/services/user-management/core"
)

type mockTokenProvider struct {
	Token string
	Err   error
}

func (m *mockTokenProvider) RequestTokenFromClientSecret(_ context.Context, clientID, clientSecret string) (string, error) {
	if m.Err != nil {
		return "", m.Err
	}
	return m.Token, nil
}

// Test authentication with valid credentials
func TestAuthenticate_Success(t *testing.T) {
	provider := &mockTokenProvider{Token: "dummy.jwt.token"}

	service := core.NewAuthService(provider)

	clientID, token, err := service.Authenticate(context.Background(), "abc.def")

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if clientID != "abc" {
		t.Errorf("expected clientID 'abc', got: %s", clientID)
	}
	if token != "dummy.jwt.token" {
		t.Errorf("expected token, got: %s", token)
	}
}

func TestAuthenticate_InvalidFormat(t *testing.T) {
	provider := &mockTokenProvider{}

	service := core.NewAuthService(provider)

	_, _, err := service.Authenticate(context.Background(), "invalidformat")
	if err == nil {
		t.Fatal("expected error for invalid format, got nil")
	}
}

func TestAuthenticate_TokenProviderFails(t *testing.T) {
	provider := &mockTokenProvider{Err: errors.New("auth0 down")}

	service := core.NewAuthService(provider)

	clientID, token, err := service.Authenticate(context.Background(), "abc.def")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if clientID != "abc" {
		t.Errorf("expected clientID 'abc', got: %s", clientID)
	}
	if token != "" {
		t.Errorf("expected empty token, got: %s", token)
	}
}

func TestIsAdminSecret(t *testing.T) {
	const correctSecret = "my-secret"
	const wrongSecret = "not-the-secret"

	hashed := sha256.Sum256([]byte(correctSecret))
	expectedHash := hex.EncodeToString(hashed[:])

	if !core.IsAdminSecret(correctSecret, expectedHash) {
		t.Error("expected IsAdminSecret to return true for correct secret")
	}

	if core.IsAdminSecret(wrongSecret, expectedHash) {
		t.Error("expected IsAdminSecret to return false for wrong secret")
	}
}
