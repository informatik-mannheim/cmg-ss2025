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
