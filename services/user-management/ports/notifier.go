package ports

type Notifier interface {
	UserRegistered(id string, role string)
	UserLoggedIn(id string)
	Event(message string)
}
