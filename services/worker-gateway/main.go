package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/informatik-mannheim/cmg-ss2025/pkg/logging"

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

	// init service and handler
	registryClient := client_http.NewRegistryClient("http://registry:8080")
	jobClient := client_http.NewJobClient("http://job:8080")
	service := core.NewWorkerGatewayService(registryClient, jobClient)
	handler := handler_http.NewHandler(service)

	// Router (mux)
	mux := http.NewServeMux()
	mux.HandleFunc("/worker/heartbeat", handler.HeartbeatHandler)
	mux.HandleFunc("/result", handler.SubmitResultHandler)
	mux.HandleFunc("/register", handler.RegisterWorkerHandler)

	// Server-Setup
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
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
