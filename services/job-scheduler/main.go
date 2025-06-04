package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	carbonintensity "github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/adapters/carbon-intensity"
	handler_http "github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/adapters/handler-http"
	interval_runner "github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/adapters/interval-runner"
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
	Port                       string
	// TODO: Add address for UserManagement; Not relevant for now, comes with phase 3
}

func main() {
	// Read environment variables
	envs, err := loadEnvVariables()
	if err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
		return
	}

	// Initialize adapters and service
	var jobAdapter ports.JobAdapter = job.NewJobAdapter(envs.JobServiceUrl)
	var workerAdapter ports.WorkerAdapter = worker.NewWorkerAdapter(envs.WorkerRegestryUrl)
	var carbonIntensityAdapter ports.CarbonIntensityAdapter = carbonintensity.NewCarbonIntensityAdapter(envs.CarbonIntensityProviderUrl)
	var service ports.JobScheduler = core.NewJobSchedulerService(
		jobAdapter,
		workerAdapter,
		carbonIntensityAdapter,
	)

	// Start the HTTP server
	srv := &http.Server{Addr: ":" + envs.Port}

	handler := handler_http.NewHandler(service)
	http.Handle("/", handler)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Print("The service is shutting down...")
		srv.Shutdown(context.Background())
		cancel() // cancel the context to stop the scheduler
	}()

	// Start the HTTP server
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()
	log.Printf("Job Scheduler is running on port %s...\n", envs.Port)

	// Start the job scheduler runner
	runner := interval_runner.NewIntervalRunner(ctx, envs.Interval, envs.Port)
	runner.RunScheduleJob()

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

	port := utils.LoadEnvOrDefault("PORT", "8080")
	portInt, err := strconv.Atoi(port)
	if err != nil || !utils.IsPortValid(portInt) {
		return envs, err
	}
	envs.Port = port

	return envs, nil
}
