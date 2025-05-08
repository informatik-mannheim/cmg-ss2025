package main

import (
	"log"
	"strconv"
	"time"

	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/adapters/notifier"
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/core"
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/ports"
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/utils"
)

func main() {
	envs, err := loadEnvVariables()
	if err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
	}
	var interval time.Duration = time.Duration(envs.Interval) // Interval in seconds

	log.Printf("Job Scheduler starting with a %d second interval...\n", envs.Interval)

	var notifier ports.Notifier = notifier.NewHttpNotifier()
	var service ports.JobScheduler = core.NewJobSchedulerService(notifier)

	ticker := time.NewTicker(interval * time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C
		service.ScheduleJob()
	}
}

func loadEnvVariables() (ports.Environments, error) {
	var envs ports.Environments = ports.Environments{}

	interval := utils.LoadEnvOrDefault("JOB_SCHEDULER_INTERVAL", "5") // Default to 5 seconds
	intervalInt, err := strconv.Atoi(interval)
	if err != nil {
		return envs, err
	}
	envs.Interval = intervalInt

	return envs, nil
}
