package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/informatik-mannheim/cmg-ss2025/pkg/logging"

	client "github.com/informatik-mannheim/cmg-ss2025/services/worker-registry/adapters/clients"
	handler "github.com/informatik-mannheim/cmg-ss2025/services/worker-registry/adapters/handler-http"
	repo "github.com/informatik-mannheim/cmg-ss2025/services/worker-registry/adapters/repo"
	inMemoryRepo "github.com/informatik-mannheim/cmg-ss2025/services/worker-registry/adapters/repo-in-memory"
	"github.com/informatik-mannheim/cmg-ss2025/services/worker-registry/core"
)

func main() {
	logging.Init("worker-registry")
	logging.Debug("Starting Worker Registry")
	zoneClient := client.NewZoneClient(os.Getenv("CARBON_INTENSITY_PROVIDER"))

	var service *core.WorkerRegistryService

	repo, err := repo.NewRepo(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		context.Background(),
	)

	if err != nil {
		log.Printf("ERROR: Failed to connect to Postgres: %v", err)
		log.Print("Could not connect to database, using in-memory fallback")
		inMemoryRepo := inMemoryRepo.NewRepo()
		service = core.NewWorkerRegistryService(inMemoryRepo, zoneClient)
	} else {
		service = core.NewWorkerRegistryService(repo, zoneClient)
	}

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
