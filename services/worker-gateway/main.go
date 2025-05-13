package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	handler_http "github.com/informatik-mannheim/cmg-ss2025/services/worker-gateway/adapters/handler-http"
	"github.com/informatik-mannheim/cmg-ss2025/services/worker-gateway/core"
	"github.com/informatik-mannheim/cmg-ss2025/services/worker-gateway/ports"
)

type testNotifier struct{}

func (t *testNotifier) UpdateWorkerStatus(ctx context.Context, req ports.HeartbeatRequest) error {
	log.Printf("[MOCK] UpdateWorkerStatus: %s -> %s", req.WorkerID, req.Status)
	return nil
}

func (t *testNotifier) FetchScheduledJobs(ctx context.Context) ([]ports.Job, error) {
	log.Println("[MOCK] FetchScheduledJobs called")
	return []ports.Job{
		{ID: "job1", Status: "SCHEDULED"},
		{ID: "job2", Status: "SCHEDULED"},
	}, nil
}

func (t *testNotifier) UpdateJob(ctx context.Context, req ports.ResultRequest) error {
	log.Printf("[MOCK] UpdateJob: %s -> %s (%s)", req.JobID, req.Status, req.Result)
	return nil
}

func (t *testNotifier) RegisterWorker(ctx context.Context, req ports.RegisterRequest) error {
	log.Printf("[MOCK] RegisterWorker: %s (%s, %s)", req.ID, req.Key, req.Location)
	return nil
}

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// init service and handler
	notifier := &testNotifier{} // TODO: replace
	service := core.NewWorkerGatewayService(notifier)
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

		log.Println("The service is shutting down...")
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Fatalf("Shutdown failed: %v", err)
		}
	}()

	log.Printf("Worker Gateway listening on port %s", port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
}
