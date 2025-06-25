package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	// "github.com/informatik-mannheim/cmg-ss2025/pkg/logging"
	handler_http "github.com/informatik-mannheim/cmg-ss2025/services/job/adapters/handler-http"
	repo_database "github.com/informatik-mannheim/cmg-ss2025/services/job/adapters/repo-database"
	repo_in_memory "github.com/informatik-mannheim/cmg-ss2025/services/job/adapters/repo-in-memory"
	"github.com/informatik-mannheim/cmg-ss2025/services/job/core"
	"github.com/informatik-mannheim/cmg-ss2025/services/job/ports"
)

func main() {
	// Initialize the job storage based on the environment variable JOB_REPO_TYPE
	var storage ports.JobStorage
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
			dbStorage, err := repo_database.NewJobStorage(host, port, user, password, dbName, sslMode, ctx)
			if err != nil {
				//logging.Warn("Postgres not ready, retrying...")
				log.Println("WARN: Postgres not ready, retrying...")
				time.Sleep(3 * time.Second)
			} else {
				//logging.Debug("Successfully connected to the PostgreSQL database")
				log.Println("DEBUG: Successfully connected to the PostgreSQL database")
				storage = dbStorage
				break
			}
		}
		if storage == nil {
			//logging.Error("Failed to connect to Postgres after multiple attempts, falling back to in-memory storage")
			log.Println("ERROR: Failed to connect to Postgres after multiple attempts, falling back to in-memory storage")
			storage = repo_in_memory.NewMockJobStorage()
		}
	} else {
		//logging.Debug("Using in-memory job storage")
		log.Println("DEBUG: Using in-memory job storage")
		// Default to in-memory storage if JOB_REPO_TYPE is not set or is not
		storage = repo_in_memory.NewMockJobStorage()
	}
	jobService, err := core.NewJobService(storage)
	if err != nil {
		log.Fatalf("could not initialize job service: %v", err)
	}

	// Set up the HTTP server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}
	handler := handler_http.NewHandler(jobService)
	if handler == nil {
		log.Fatal("could not create handler")
	}
	server := &http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}
	// Start the server in a goroutine so that it doesn't block the main thread
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Print("The service is shutting down...")
		server.Shutdown(context.Background())
	}()

	log.Print("listening...")
	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("could not listen on %s: %v\n", server.Addr, err)
	}
	log.Print("Done")
}
