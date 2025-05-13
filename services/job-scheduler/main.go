package main

import (
	"log"
	"strconv"
	"time"

	carbonintensity "github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/adapters/carbon-intensity"
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/adapters/job"
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/adapters/notifier"
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/adapters/worker"
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/core"
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/model"
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/ports"
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/utils"
)

func main() {
	envs, err := loadEnvVariables()
	if err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
		return
	}
	var interval time.Duration = time.Duration(envs.Interval) // Interval in seconds

	log.Printf("Job Scheduler starting with a %d second interval...\n", envs.Interval)

	var notifier ports.Notifier = notifier.NewNotifier()
	var jobAdapter ports.JobAdapter = job.NewJobAdapter(envs)
	var workerAdapter ports.WorkerAdapter = worker.NewWorkerAdapter(envs)
	var carbonIntensityAdapter ports.CarbonIntensityAdapter = carbonintensity.NewCarbonIntensityAdapter(envs)
	var service ports.JobScheduler = core.NewJobSchedulerService(
		jobAdapter,
		workerAdapter,
		carbonIntensityAdapter,
		notifier,
	)

	ticker := time.NewTicker(interval * time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C
		service.ScheduleJob()
	}
}

func loadEnvVariables() (model.Environments, error) {
	var envs model.Environments = model.Environments{}

	interval := utils.LoadEnvOrDefault("JOB_SCHEDULER_INTERVAL", "5") // Default to 5 seconds
	intervalInt, err := strconv.Atoi(interval)
	if err != nil {
		return envs, err
	}
	envs.Interval = intervalInt

	workerRegistry, err := utils.LoadEnvRequired("WORKER_REGISTRY")
	if err != nil || !utils.IsUrlValid(workerRegistry) {
		return envs, err
	}
	envs.WorkerRegestryUrl = workerRegistry

	jobService, err := utils.LoadEnvRequired("JOB_SERVICE")
	if err != nil || !utils.IsUrlValid(jobService) {
		return envs, err
	}
	envs.JobServiceUrl = jobService

	carbonProvider, err := utils.LoadEnvRequired("CARBON_INTENSITY_PROVIDER")
	if err != nil || !utils.IsUrlValid(carbonProvider) {
		return envs, err
	}
	envs.CarbonIntensityProviderUrl = carbonProvider

	return envs, nil
}
