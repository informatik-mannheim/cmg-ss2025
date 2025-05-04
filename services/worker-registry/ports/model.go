package ports

type Worker struct {
	Id     string `json:"id"`
	Status string `json:"status"`
	Zone   string `json:"zone"`
}
