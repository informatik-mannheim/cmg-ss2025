package ports

// Has no real use outside of testing
type Consumer struct {
	Id         string
	IntProp    int
	StringProp string
}
type JobStatus string

type ContainerImage struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}
