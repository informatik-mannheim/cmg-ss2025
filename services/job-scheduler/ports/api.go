package ports

// api.go - For possible future interaction (like metrics or CLI)

type Api interface {
	ScheduleJob() error
}
