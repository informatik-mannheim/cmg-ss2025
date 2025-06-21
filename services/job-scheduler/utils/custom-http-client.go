package utils

import (
	"net/http"

	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/ports"
)

type AuthTransport struct {
	Base    http.RoundTripper
	Adapter ports.AuthAdapter
}

var _ http.RoundTripper = (*AuthTransport)(nil)

func NewAuthTransport(base http.RoundTripper, adapter ports.AuthAdapter) *AuthTransport {
	return &AuthTransport{
		Base:    base,
		Adapter: adapter,
	}
}

func GetCustomHttpClient(adapter ports.AuthAdapter) *http.Client {
	return &http.Client{
		Transport: NewAuthTransport(nil, adapter),
	}
}

func (t *AuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req2 := req.Clone(req.Context())
	req2.AddCookie(&http.Cookie{Name: "token", Value: t.Adapter.GetToken()})

	resp, err := t.base().RoundTrip(req2)
	if err != nil || resp.StatusCode != http.StatusUnauthorized {
		return resp, err
	}

	if err := t.Adapter.Authenticate(); err != nil {
		return resp, err
	}

	req3 := req.Clone(req.Context())
	req3.AddCookie(&http.Cookie{Name: "token", Value: t.Adapter.GetToken()})
	return t.base().RoundTrip(req3)
}

func (t *AuthTransport) base() http.RoundTripper {
	if t.Base != nil {
		return t.Base
	}
	return http.DefaultTransport
}
