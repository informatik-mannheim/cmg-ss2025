package ports

// FIXME: add errors

type JobScheduler interface {
	ScheduleJob() error
}
