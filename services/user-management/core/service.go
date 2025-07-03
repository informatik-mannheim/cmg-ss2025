package core

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"strings"

	"github.com/informatik-mannheim/cmg-ss2025/services/user-management/ports"
)

type AuthService struct {
	TokenProvider ports.TokenProvider
}

func NewAuthService(tokenProvider ports.TokenProvider) *AuthService {
	return &AuthService{
		TokenProvider: tokenProvider,
	}
}

func (s *AuthService) Authenticate(ctx context.Context, credentials string) (string, string, error) {
	clientID, clientSecret, err := splitCredentials(credentials)
	if err != nil {
		return "", "", err
	}

	token, err := s.TokenProvider.RequestTokenFromClientSecret(ctx, clientID, clientSecret)
	if err != nil {
		return clientID, "", err
	}

	return clientID, token, nil
}

func splitCredentials(secret string) (string, string, error) {
	parts := strings.SplitN(secret, ".", 2)
	if len(parts) != 2 {
		return "", "", errors.New("invalid credentials format")
	}
	return parts[0], parts[1], nil
}

func IsAdminSecret(input, expectedHash string) bool {
	hash := sha256.Sum256([]byte(input))
	actualHash := hex.EncodeToString(hash[:])
	return subtle.ConstantTimeCompare([]byte(actualHash), []byte(expectedHash)) == 1
}
