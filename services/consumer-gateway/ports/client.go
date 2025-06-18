package ports

import "context"

type ZoneClient interface {
	GetZone(ctx context.Context, req ZoneRequest) (ZoneResponse, error)
}

type LoginClientRequest struct {
	Secret string `json:"secret"`
}

type LoginClientResponse struct {
	Token string `json:"token"`
}

type LoginClient interface {
	Login(ctx context.Context, req LoginClientRequest) (LoginClientResponse, error)
}

type JobClient interface {
	GetJobOutcome(ctx context.Context, jobID string) (JobOutcomeResponse, error)
	CreateJob(ctx context.Context, req CreateJobRequest) (CreateJobResponse, error)
}
