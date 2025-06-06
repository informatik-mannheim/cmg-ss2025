package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	client "github.com/informatik-mannheim/cmg-ss2025/services/worker-registry/adapters/clients"
	handler "github.com/informatik-mannheim/cmg-ss2025/services/worker-registry/adapters/handler-http"
	notifier "github.com/informatik-mannheim/cmg-ss2025/services/worker-registry/adapters/notifier"
	repo "github.com/informatik-mannheim/cmg-ss2025/services/worker-registry/adapters/repo"
	inMemoryRepo "github.com/informatik-mannheim/cmg-ss2025/services/worker-registry/adapters/repo-in-memory"
	"github.com/informatik-mannheim/cmg-ss2025/services/worker-registry/core"
)

func main() {
	notifier := notifier.NewNotifier()
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
		service = core.NewWorkerRegistryService(inMemoryRepo, notifier, zoneClient)
	} else {
		service = core.NewWorkerRegistryService(repo, notifier, zoneClient)
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

		log.Print("The service is shutting down...")
		srv.Shutdown(context.Background())
	}()

	log.Print("listening...")
	srv.ListenAndServe()
	log.Print("Done")
}
