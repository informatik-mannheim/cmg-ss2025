package ports

import (
	"context"
)

type UserService interface {
	GetToken(ctx context.Context, req GetTokenRequest) (GetTokenResponse, error)
}

type GetTokenRequest struct {
	Secret string `json:"secret"`
}

type GetTokenResponse struct {
	Token string `json:"token"`
}
