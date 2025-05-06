package notifier

import "log"

// Notifier is a basic interface for sending user notifications.
type Notifier interface {
	UserRegistered(id string, role string)
	UserLoggedIn(id string)
	Event(message string)
}

// StdoutNotifier is a simple implementation that logs to stdout.
type StdoutNotifier struct{}

// UserRegistered logs the registration event to stdout.
func (n *StdoutNotifier) UserRegistered(id string, role string) {
	log.Printf("[Notifier] New user registered: ID=%s, Role=%s", id, role)
}

func (n *StdoutNotifier) UserLoggedIn(id string) {
	log.Printf("[Notifier] User logged in: ID=%s", id)
}

// New returns a default notifier implementation.
func New() Notifier {
	return &StdoutNotifier{}
}

func (n *StdoutNotifier) Event(message string) {
	log.Printf("[Notifier] %s", message)
}
