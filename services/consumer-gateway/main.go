package main

import (
	"context"
	"fmt"
	"github.com/informatik-mannheim/cmg-ss2025/pkg/auth"
	"github.com/informatik-mannheim/cmg-ss2025/pkg/logging"

	"github.com/informatik-mannheim/cmg-ss2025/pkg/tracing/tracing"
	jobclient "github.com/informatik-mannheim/cmg-ss2025/services/consumer-gateway/adapters/client-http"
	handler_http "github.com/informatik-mannheim/cmg-ss2025/services/consumer-gateway/adapters/handler-http"
	"github.com/informatik-mannheim/cmg-ss2025/services/consumer-gateway/core"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func secure(h http.HandlerFunc) http.Handler {
	return auth.AuthMiddleware(h)
}

func main() {

	logging.Init("consumer-gateway")

	jaeger := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if jaeger == "" {
		logging.Warn("Environment variable OTEL_EXPORTER_OTLP_ENDPOINT is not set")
	}

	shutdown, err := tracing.Init("test-service", jaeger) // Init tracer
	if err != nil {
		logging.Warn("Tracing init failed:", err)
	}
	defer shutdown(context.Background())

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello from Jaeger-traced service!")
	})

	log.Println("Listening on :8080...")
	http.ListenAndServe(":8080", tracing.Middleware(mux)) // apply tracing middleware

	jwksUrl := os.Getenv("JWKS_URL")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	err = auth.InitJWKS(jwksUrl)
	if err != nil {
		logging.Error("Failed to initialize JWKS: %v", err)
	}

	// set port manually with "export PORT"
	job := jobclient.NewJobClient(os.Getenv("JOB_SERVICE"))
	user := jobclient.NewLoginClient(os.Getenv("USER_MANAGEMENT_SERVICE"))
	zone := jobclient.NewZoneClient(os.Getenv("CARBON_INTENSITY_PROVIDER"))
	service := core.NewConsumerService(job, zone, user)
	handler := handler_http.NewHandler(service)

	mux = http.NewServeMux()
	mux.Handle("/", handler)

	srv := &http.Server{Addr: ":" + port, Handler: mux}

	mux.Handle("/jobs", secure(handler.HandleCreateJobRequest))
	mux.Handle("jobs/{id}/result", secure(handler.HandleGetJobOutcomeRequest))
	mux.HandleFunc("/auth/login", handler.HandleLoginRequest)

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		logging.Debug("The service is shutting down...")
		err := srv.Shutdown(context.Background())
		if err != nil {
			return
		}
	}()

	logging.Debug("listening on port " + port + " ...")

	err = srv.ListenAndServe()
	if err != nil {
		return
	}
	logging.Debug("Done")
}
