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

func secure(h http.HandlerFunc) http.Handler {
	return auth.AuthMiddleware(h)
}

func main() {

	jwksUrl := os.Getenv("JWKS_URL")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	err := auth.InitJWKS(jwksUrl)
	if err != nil {
		log.Fatalf("Failed to initialize JWKS: %v", err)
	}

	// set port manually with "export PORT"
	job := jobclient.NewJobClient(os.Getenv("JOB_SERVICE"))
	user := jobclient.NewLoginClient(os.Getenv("USER_MANAGEMENT_SERVICE"))
	zone := jobclient.NewZoneClient(os.Getenv("CARBON_INTENSITY_PROVIDER"))
	service := core.NewConsumerService(job, zone, user)
	handler := handler_http.NewHandler(service)

	mux := http.NewServeMux()
	mux.Handle("/", handler)

	srv := &http.Server{Addr: ":" + port, Handler: mux}

	mux.Handle("/jobs", secure(handler.HandleCreateJobRequest))
	mux.Handle("jobs/{id}/result", secure(handler.HandleGetJobOutcomeRequest))
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
