package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/informatik-mannheim/cmg-ss2025/pkg/logging"
	"github.com/informatik-mannheim/cmg-ss2025/pkg/tracing/tracing"
	handler_http "github.com/informatik-mannheim/cmg-ss2025/services/job/adapters/handler-http"
	repo_database "github.com/informatik-mannheim/cmg-ss2025/services/job/adapters/repo-database"
	repo_in_memory "github.com/informatik-mannheim/cmg-ss2025/services/job/adapters/repo-in-memory"
	"github.com/informatik-mannheim/cmg-ss2025/services/job/core"
	"github.com/informatik-mannheim/cmg-ss2025/services/job/ports"
)

func main() {
	// Initialize Logging
	logging.Init("job-service")
	logging.Debug("starting job-service")

	// Initialize Tracing
	jaeger := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if jaeger == "" {
		logging.Error("Environment variable OTEL_EXPORTER_OTLP_ENDPOINT is not set")
	}
	shutdown, err := tracing.Init("job-service", jaeger)
	if err != nil {
		errorMessage := "could not initialize tracing: " + err.Error()
		logging.Error(errorMessage)
	}
	defer shutdown(context.Background())

	// Initialize the job storage based on the environment variable JOB_REPO_TYPE
	var storage ports.JobStorage
	var dbError error
	repoType := os.Getenv("JOB_REPO_TYPE")
	if repoType == "postgres" {
		host := os.Getenv("DB_HOST")
		port := os.Getenv("DB_PORT")
		user := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASSWORD")
		dbName := os.Getenv("DB_NAME")
		sslMode := os.Getenv("SSL_MODE")
		ctx := context.Background()

		maxRetries := 10
		for i := 0; i < maxRetries; i++ {
			var dbStorage ports.JobStorage
			dbStorage, dbError = repo_database.NewJobStorage(host, port, user, password, dbName, sslMode, ctx)
			if dbError != nil {
				logging.Warn("Postgres not ready, retrying...")
				time.Sleep(3 * time.Second)
			} else {
				logging.Debug("Successfully connected to the PostgreSQL database")
				storage = dbStorage
				break
			}
		}
		if storage == nil && dbError != nil {
			errorMessage := "could not connect to Postgres after multiple attempts: " + dbError.Error()
			logging.Error(errorMessage)
			// Fallback to in-memory storage if Postgres connection fails
			logging.Warn("Falling back to in-memory job storage")
			storage = repo_in_memory.NewMockJobStorage()
		}
	} else {
		logging.Debug("Using in-memory job storage")
		storage = repo_in_memory.NewMockJobStorage()
	}
	jobService, err := core.NewJobService(storage)
	if err != nil {
		errorMessage := "could not initialize job service: " + err.Error()
		logging.Error(errorMessage)
	}

	// Set up the HTTP server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}
	handler := handler_http.NewHandler(jobService)
	if handler == nil {
		logging.Warn("could not create handler, service is shutting down")
		os.Exit(1)
	}
	tracingHandler := tracing.Middleware(handler)
	server := &http.Server{
		Addr:    ":" + port,
		Handler: tracingHandler,
	}

	// Start the server in a goroutine so that it doesn't block the main thread
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		logging.Debug("The service is shutting down...")
		server.Shutdown(context.Background())
	}()

	logging.Debug("listening...")
	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		errorMessage := "could not listen on " + server.Addr + ": " + err.Error()
		logging.Error(errorMessage)
	}
	logging.Debug("Done")
}
