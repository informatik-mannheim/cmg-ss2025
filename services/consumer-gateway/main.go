package main

import (
	"context"
	"github.com/informatik-mannheim/cmg-ss2025/pkg/auth"
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
	err := auth.InitJWKS("https://dev-jqhwcu7xuwgdqi56.eu.auth0.com/.well-known/jwks.json")
	if err != nil {
		log.Fatalf("Failed to initialize JWKS: %v", err)
	}

	// set port manually with "export PORT"
	job := jobclient.NewJobClient(os.Getenv("JOB_SERVICE") + port)
	user := jobclient.NewLoginClient(os.Getenv("USER_MANAGEMENT_SERVICE") + port)
	zone := jobclient.NewZoneClient(os.Getenv("CARBON_INTENSITY_PROVIDER") + port)
	service := core.NewConsumerService(job, zone, user)
	handler := handler_http.NewHandler(service, service, service)

	mux := http.NewServeMux()
	mux.Handle("/", handler)

	srv := &http.Server{Addr: ":" + port, Handler: mux}

	mux.Handle("/jobs", auth.AuthMiddleware(http.HandlerFunc(handler.HandleCreateJobRequest)))
	mux.Handle("jobs/{id}/result", auth.AuthMiddleware(http.HandlerFunc(handler.HandleGetJobOutcomeRequest)))
	mux.HandleFunc("/auth/login", handler.HandleLoginRequest)

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Print("The service is shutting down...")
		err := srv.Shutdown(context.Background())
		if err != nil {
			return
		}
	}()


	log.Print("listening on port " + port + " ...")

	err = srv.ListenAndServe()
	if err != nil {
		return
	}
	log.Print("Done")
}
