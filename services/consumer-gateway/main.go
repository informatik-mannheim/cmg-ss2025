package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	jobclient "github.com/informatik-mannheim/cmg-ss2025/services/consumer-gateway/adapters/client-http"
	handler_http "github.com/informatik-mannheim/cmg-ss2025/services/consumer-gateway/adapters/handler-http"
	"github.com/informatik-mannheim/cmg-ss2025/services/consumer-gateway/core"
)

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	job := jobclient.NewJobClient("http://job:8080")
	user := jobclient.NewLoginClient("http://auth/login:8080")
	zone := jobclient.NewZoneClient("http://carbon-intensity-provider:8080")
	service := core.NewConsumerService(job, zone, user)
	handler := handler_http.NewHandler(service, service, service)

	// Router (mux)
	mux := http.NewServeMux()
	mux.HandleFunc("/jobs", handler.HandleCreateJobRequest)
	mux.HandleFunc("jobs/{id}/result", handler.HandleGetJobOutcomeRequest)
	mux.HandleFunc("/auth/login", handler.HandleLoginRequest)

	srv := &http.Server{Addr: ":8080"}

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
