package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/informatik-mannheim/cmg-ss2025/pkg/auth"
	"github.com/informatik-mannheim/cmg-ss2025/pkg/logging"
	"github.com/informatik-mannheim/cmg-ss2025/pkg/tracing/tracing"

	client_http "github.com/informatik-mannheim/cmg-ss2025/services/worker-gateway/adapters/client-http"
	handler_http "github.com/informatik-mannheim/cmg-ss2025/services/worker-gateway/adapters/handler-http"
	"github.com/informatik-mannheim/cmg-ss2025/services/worker-gateway/core"
)

func main() {
	logging.Init("worker-gateway")
	logging.Debug("Service started")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	jaeger := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if jaeger == "" {
		logging.Error("Environment variable OTEL_EXPORTER_OTLP_ENDPOINT is not set")
	}

	shutdown, err := tracing.Init("worker-gateway", jaeger)
	if err != nil {
		logging.Error("Tracing init failed:", err)
	}
	defer shutdown(context.Background())

	jwksUrl := os.Getenv("JWKS_URL")

	err = auth.InitJWKS(jwksUrl)

	if err != nil {
		logging.Error("Failed to initialize JWKS: " + err.Error())
		return
	}

	// init service and handler
	registryClient := client_http.NewRegistryClient(os.Getenv("WORKER_REGISTRY"))
	jobClient := client_http.NewJobClient(os.Getenv("JOB_SERVICE"))
	userClient := client_http.NewUserClient(os.Getenv("USER_MANAGEMENT_SERVICE"))
	service := core.NewWorkerGatewayService(registryClient, jobClient, userClient)
	handler := handler_http.NewHandler(service)

	// Router (mux)
	mux := http.NewServeMux()
	mux.Handle("/worker/heartbeat", auth.AuthMiddleware(http.HandlerFunc(handler.HeartbeatHandler)))
	mux.Handle("/result", auth.AuthMiddleware(http.HandlerFunc(handler.SubmitResultHandler)))
	mux.Handle("/register", http.HandlerFunc(handler.RegisterWorkerHandler))

	// Wrap router with tracing middleware
	tracingHandler := tracing.Middleware(mux)

	// Server-Setup
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: tracingHandler,
	}

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		logging.Debug("The service is shutting down...")
		if err := srv.Shutdown(context.Background()); err != nil {
			logging.Error("Shutdown failed", "err", err)
		}
	}()

	logging.Debug("Worker Gateway listening", "port", port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logging.Error("Server error", "err", err)
	}
}
