package ports

import "context"

type AuthProvider interface {
	RequestTokenFromCredentials(ctx context.Context, credentials string) (string, error)
}
