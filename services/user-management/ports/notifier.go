package ports

import "context"

type Notifier interface {
	UserRegistered(id string, role string, ctx context.Context)
	UserLoggedIn(id string, ctx context.Context)
	Event(message string, ctx context.Context)
}
