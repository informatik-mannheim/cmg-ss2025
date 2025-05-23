package cli

type CreateJobRequest struct {
	ImageID      string            `json:"image_id"`
	JobName      string            `json:"job_name"`
	CreationZone string            `json:"creation_zone"`
	Parameters   map[string]string `json:"parameters"`
}
