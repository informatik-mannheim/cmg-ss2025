package main

import (
	"log"
	"strconv"
	"time"

	carbonintensity "github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/adapters/carbon-intensity"
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/adapters/job"
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/adapters/worker"
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/core"
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/ports"
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/utils"
)

// This struct is used to define the environment variables for the job-scheduler
type Environments struct {
	Interval                   int
	WorkerRegestryUrl          string
	JobServiceUrl              string
	CarbonIntensityProviderUrl string
	// TODO: Add address for UserManagement; Not relevant for now, comes with phase 3
}

func main() {
	envs, err := loadEnvVariables()
	if err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
		return
	}
	var interval time.Duration = time.Duration(envs.Interval) // Interval in seconds

	log.Printf("Job Scheduler starting with a %d second interval...\n", envs.Interval)

	var jobAdapter ports.JobAdapter = job.NewJobAdapter(envs.JobServiceUrl)
	var workerAdapter ports.WorkerAdapter = worker.NewWorkerAdapter(envs.WorkerRegestryUrl)
	var carbonIntensityAdapter ports.CarbonIntensityAdapter = carbonintensity.NewCarbonIntensityAdapter(envs.CarbonIntensityProviderUrl)
	var service ports.JobScheduler = core.NewJobSchedulerService(
		jobAdapter,
		workerAdapter,
		carbonIntensityAdapter,
	)

	ticker := time.NewTicker(interval * time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C
		service.ScheduleJob()
	}
}

func loadEnvVariables() (Environments, error) {
	var envs Environments = Environments{}

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
