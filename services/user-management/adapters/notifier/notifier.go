package notifier

import (
	"log"

	"github.com/informatik-mannheim/cmg-ss2025/services/user-management/ports"
)

type StdoutNotifier struct{}

func New() ports.Notifier {
	return &StdoutNotifier{}
}

func (n *StdoutNotifier) UserRegistered(id string, role string) {
	log.Printf("[Notifier] New user registered: ID=%s, Role=%s", id, role)
}

func (n *StdoutNotifier) UserLoggedIn(id string) {
	log.Printf("[Notifier] User logged in: ID=%s", id)
}

func (n *StdoutNotifier) Event(message string) {
	log.Printf("[Notifier] %s", message)
}
