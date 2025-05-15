package ports

type JobScheduler interface {
	ScheduleJob() error
}
