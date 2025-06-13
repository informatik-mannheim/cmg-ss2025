package auth

import (
	"fmt"

	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/ports"
)

func GetAuthEndpoint(base string) string {
	return fmt.Sprintf("%s/auth/login", base)
}

type AuthAdapter struct {
	baseUrl string
	secret  string
}

var _ ports.AuthAdapter = (*AuthAdapter)(nil)

func NewAuthAdapter(baseUrl, secret string) *AuthAdapter {
	return &AuthAdapter{
		baseUrl: baseUrl,
		secret:  secret,
	}
}

func (adapter *AuthAdapter) Authenticate() error {

	// FIXME: implement

	return nil
}
