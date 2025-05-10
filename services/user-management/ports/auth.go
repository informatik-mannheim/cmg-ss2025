package ports

type AuthProvider interface {
	RequestTokenFromCredentials(credentials string) (string, error)
}
