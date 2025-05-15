package ports

type Role string

const (
	Consumer     Role = "consumer"
	Provider     Role = "provider"
	JobScheduler Role = "job scheduler"
)

type UserManagement struct {
	ID     string `json:"id"`
	Role   Role   `json:"role"`
	Secret string `json:"secret"`
}

// registerRequest is the request payload for user registration.
// It contains the role of the user to be registered.
// The role can be either Consumer, Provider, or JobScheduler.
type RegisterRequest struct {
	Role Role `json:"role"`
}

// registerResponse is the response payload for user registration.
// It contains the generated user ID.
type RegisterResponse struct {
	Secret string `json:"secret"`
}

// loginRequest is the request payload for user login.
// It contains the secret used for authentication.
// The secret should be in the format "clientID.clientSecret".
type LoginRequest struct {
	Secret string `json:"secret"`
}

// loginResponse is the response payload for user login.
// It contains the token received from Auth0.
type LoginResponse struct {
	Token string `json:"token"`
}
