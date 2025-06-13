package ports

type AuthAdapter interface {
	Authenticate() error
}
