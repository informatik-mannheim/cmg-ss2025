package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	repo "github.com/informatik-mannheim/cmg-ss2025/services/carbon-intensity-provider/adapters/repo-in-memory"
	"github.com/informatik-mannheim/cmg-ss2025/services/carbon-intensity-provider/api"
	"github.com/informatik-mannheim/cmg-ss2025/services/carbon-intensity-provider/core"
)

func main() {
	repository := repo.NewRepo()
	service := core.NewCarbonIntensityService(repository)

	// Preload fixed zones manually as test(for Assignment II)
	service.AddOrUpdateZone("DE", 140.5)
	service.AddOrUpdateZone("FR", 135.2)
	service.AddOrUpdateZone("US-NY-NYIS", 128.9)

	// Start server
	router := api.NewHandler(service)
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Println("Carbon Intensity Provider running on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-stop
	log.Println("Shutdown signal received...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server exited gracefully.")
}
