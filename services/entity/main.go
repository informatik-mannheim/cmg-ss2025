package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	handler_http "github.com/informatik-mannheim/cmg-ss2025/services/entity/adapters/handler-http"
	repo "github.com/informatik-mannheim/cmg-ss2025/services/entity/adapters/repo-in-memory"
	"github.com/informatik-mannheim/cmg-ss2025/services/entity/core"
)

func main() {
	coreService := core.NewEntityService(repo.NewRepo(), nil)
	handler := handler_http.NewHandler(coreService)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	// Channel to listen for interrupt or terminate signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Run the server in a goroutine
	go func() {
		log.Println("Listening on :8080...")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for signal
	<-stop
	log.Println("Shutdown signal received...")

	// Give 5 seconds for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server exited gracefully.")
}
