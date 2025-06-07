package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/informatik-mannheim/cmg-ss2025/pkg/logging"

	client "github.com/informatik-mannheim/cmg-ss2025/services/worker-registry/adapters/clients"
	handler "github.com/informatik-mannheim/cmg-ss2025/services/worker-registry/adapters/handler-http"
	repo_pg "github.com/informatik-mannheim/cmg-ss2025/services/worker-registry/adapters/repo"
	repo_mem "github.com/informatik-mannheim/cmg-ss2025/services/worker-registry/adapters/repo-in-memory"
	"github.com/informatik-mannheim/cmg-ss2025/services/worker-registry/core"
	"github.com/informatik-mannheim/cmg-ss2025/services/worker-registry/ports"
)

func main() {
	logging.Init("worker-registry")
	logging.Debug("Starting Worker Registry")
	zoneClient := client.NewZoneClient(os.Getenv("CARBON_INTENSITY_PROVIDER"))

	var service *core.WorkerRegistryService
	var dbRepo ports.Repo

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	maxRetries := 10

	for i := 0; i < maxRetries; i++ {
		r, err := repo_pg.NewRepo(dbHost, dbPort, dbUser, dbPassword, dbName, context.Background())
		if err == nil {
			dbRepo = r
			break
		}
		logging.Warn("Postgres not ready, retrying...")
		time.Sleep(2 * time.Second)
	}

	if dbRepo == nil {
		logging.Warn("Failed to connect to Postgres")
		logging.Warn("Falling back to in-memory repository")
		dbRepo = repo_mem.NewRepo()
	}

	service = core.NewWorkerRegistryService(dbRepo, zoneClient)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	srv := &http.Server{Addr: ":" + port}

	h := handler.NewHandler(service)
	http.Handle("/", h)

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		logging.Debug("The service is shutting down...")
		srv.Shutdown(context.Background())
	}()

	logging.Debug("Listening...")
	srv.ListenAndServe()
	logging.Debug("Done")
}
