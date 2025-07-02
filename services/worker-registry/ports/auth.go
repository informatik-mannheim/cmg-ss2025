package ports

type AuthAdapter interface {
	Authenticate() error
	GetToken() string
}

type GetAuthToken struct {
	Secret string `json:"secret"`
}

type AuthTokenResponse struct {
	Token string `json:"token"`
}
