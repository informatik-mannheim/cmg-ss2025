package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	handler_http "github.com/informatik-mannheim/cmg-ss2025/services/job/adapters/handler-http"
	repo_in_memory "github.com/informatik-mannheim/cmg-ss2025/services/job/adapters/repo-in-memory"
	"github.com/informatik-mannheim/cmg-ss2025/services/job/core"
)

func main() {
	storage := repo_in_memory.NewMockJobStorage()
	jobService, err := core.NewJobService(storage)
	if err != nil {
		log.Fatalf("could not initialize job service: %v", err)
	}

	handler := handler_http.NewHandler(jobService)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Print("The service is shutting down...")
		srv.Shutdown(context.Background())
	}()

	log.Print("listening...")
	err = srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("could not listen on %s: %v\n", srv.Addr, err)
	}
	log.Print("Done")
}
