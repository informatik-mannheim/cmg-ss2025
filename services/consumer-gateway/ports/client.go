package ports

import "context"

type ZoneClient interface {
	GetZone(ctx context.Context, req ZoneRequest) (ZoneResponse, error)
}

type LoginClient interface {
	Login(ctx context.Context, req ConsumerLoginRequest) (LoginResponse, error)
}

type JobClient interface {
	GetJobOutcome(ctx context.Context, jobID string) (JobOutcomeResponse, error)
	CreateJob(ctx context.Context, req CreateJobRequest) (CreateJobResponse, error)
}
