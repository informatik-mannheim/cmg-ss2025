package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	handler "github.com/informatik-mannheim/cmg-ss2025/services/worker-registry/adapters/handler-http"
	notifier "github.com/informatik-mannheim/cmg-ss2025/services/worker-registry/adapters/notifier"
	repo "github.com/informatik-mannheim/cmg-ss2025/services/worker-registry/adapters/repo-in-memory"
	"github.com/informatik-mannheim/cmg-ss2025/services/worker-registry/core"
	"github.com/informatik-mannheim/cmg-ss2025/services/worker-registry/ports"
)

func main() {
	repository := repo.NewRepo()
	notifier := notifier.NewHttpNotifier()
	service := core.NewWorkerRegistryService(repository, notifier)

	CreateDummyWorkers(*service)

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

// Preload fixed workers manually as test(for Assignment II)
func CreateDummyWorkers(service core.WorkerRegistryService) {
	service.CreateWorker("DE", context.Background())
	service.CreateWorker("EN", context.Background())
	service.CreateWorker("DE", context.Background())
	service.CreateWorker("DE", context.Background())
	service.UpdateWorkerStatus("3", ports.StatusRunning, context.Background())
}
