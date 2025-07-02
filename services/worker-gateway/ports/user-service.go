package ports

import (
	"context"
)

type UserService interface {
	GetToken(ctx context.Context) (string, error)
}

type GetTokenRequest struct {
	Role string `json:"role"`
}

type GetTokenResponse struct {
	Secret string `json:"secret"`
}
