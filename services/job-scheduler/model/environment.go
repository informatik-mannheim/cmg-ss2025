package model

// This struct is used to define the environment variables for the job-scheduler
type Environments struct {
	Interval                   int
	WorkerRegestryUrl          string
	JobServiceUrl              string
	CarbonIntensityProviderUrl string
	// TODO: Add address for UserManagement; Not relevant for now, comes with phase 3
}
