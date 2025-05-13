package ports

import "context"

type TokenProvider interface {
	RequestTokenFromClientSecret(ctx context.Context, clientID, clientSecret string) (string, error)
}
