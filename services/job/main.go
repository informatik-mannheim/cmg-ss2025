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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}

	handler := handler_http.NewHandler(jobService)
	if handler == nil {
		log.Fatal("could not create handler")
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Print("The service is shutting down...")
		server.Shutdown(context.Background())
	}()

	log.Print("listening...")
	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("could not listen on %s: %v\n", server.Addr, err)
	}
	log.Print("Done")
}
