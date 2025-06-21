package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/informatik-mannheim/cmg-ss2025/pkg/logging"
	"github.com/informatik-mannheim/cmg-ss2025/pkg/tracing/tracing"

	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/adapters/auth"
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
	UserManagementUrl          string
	AuthToken                  string
	OTLPExporterOtlpEndpoint   string // OpenTelemetry endpoint for tracing
}

func main() {
	logging.Init("job-scheduler")
	logging.Debug("Job Scheduler service is starting...")

	// Read environment variables
	envs, err := loadEnvVariables()
	if err != nil {
		logging.Error("Failed to load environment variables: %v", err)
		return
	}

	shutdown, err := tracing.Init("job-scheduler", envs.OTLPExporterOtlpEndpoint)
	if err != nil {
		logging.Error("Failed to initialize tracing: %v", err)
		return
	}
	defer shutdown(context.Background())

	// Initialize adapters and service
	var authAdapter ports.AuthAdapter = auth.NewAuthAdapter(envs.UserManagementUrl, envs.AuthToken)
	var customClient http.Client = *utils.GetCustomHttpClient(authAdapter)

	var jobAdapter ports.JobAdapter = job.NewJobAdapter(customClient, envs.JobServiceUrl)
	var workerAdapter ports.WorkerAdapter = worker.NewWorkerAdapter(customClient, envs.WorkerRegestryUrl)
	var carbonIntensityAdapter ports.CarbonIntensityAdapter = carbonintensity.NewCarbonIntensityAdapter(customClient, envs.CarbonIntensityProviderUrl)
	var service ports.JobScheduler = core.NewJobSchedulerService(
		jobAdapter,
		workerAdapter,
		carbonIntensityAdapter,
	)

	// Start the HTTP server
	srv := &http.Server{Addr: ":" + envs.Port}

	handler := handler_http.NewHandler(service)
	tracingHandler := tracing.Middleware(handler)
	http.Handle("/", tracingHandler)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		logging.Debug("Received shutdown signal, shutting down the server...")
		srv.Shutdown(context.Background())
		cancel() // cancel the context to stop the scheduler
	}()

	// Start the HTTP server
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logging.Error("HTTP server error: %v", err)
		}
	}()
	logging.Debug("Job Scheduler is running on port %s...", envs.Port)

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

	userManagementUrl, err := utils.LoadEnvRequired("USER_MANAGEMENT_URL")
	if err != nil || !utils.IsUrlValid(userManagementUrl) {
		return envs, err
	}
	envs.UserManagementUrl = userManagementUrl

	authToken, err := utils.LoadEnvRequired("AUTH_TOKEN")
	if err != nil || authToken == "" {
		return envs, err
	}
	envs.AuthToken = authToken

	otlpEndpoint, err := utils.LoadEnvRequired("OTEL_EXPORTER_OTLP_ENDPOINT")
	if err != nil || !utils.IsUrlValid(otlpEndpoint) {
		return envs, err
	}
	envs.OTLPExporterOtlpEndpoint = otlpEndpoint

	return envs, nil
}
