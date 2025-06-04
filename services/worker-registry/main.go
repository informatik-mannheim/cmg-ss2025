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
	repo "github.com/informatik-mannheim/cmg-ss2025/services/worker-registry/adapters/repo-in-memory"
	"github.com/informatik-mannheim/cmg-ss2025/services/worker-registry/core"
)

func main() {
	repository := repo.NewRepo()
	notifier := notifier.NewNotifier()
	zoneClient := client.NewZoneClient(os.Getenv("CARBON_INTENSITY_PROVIDER"))
	service := core.NewWorkerRegistryService(repository, notifier, zoneClient)

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
