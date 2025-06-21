package auth

import (
	"fmt"
	"net/http"

	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/ports"
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/utils"
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

	result, _, err := utils.PostRequest[ports.GetAuthToken, ports.AuthTokenResponse](&http.Client{}, endpoint, data)
	if err != nil {
		return fmt.Errorf("failed to authenticate: %w", err)
	}

	adapter.token = result.Token

	return nil
}

func (adapter *AuthAdapter) GetToken() string {
	return adapter.token
}
