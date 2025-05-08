package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	handler_http "github.com/informatik-mannheim/cmg-ss2025/services/worker-gateway/adapters/handler-http"
	repo "github.com/informatik-mannheim/cmg-ss2025/services/worker-gateway/adapters/repo-in-memory"
	"github.com/informatik-mannheim/cmg-ss2025/services/worker-gateway/core"
)

func main() {

	core := core.NewWorkerGatewayService(repo.NewRepo(), nil)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}
	srv := &http.Server{Addr: ":" + port}

	h := handler_http.NewHandler(core)
	http.Handle("/", h)

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Print("The service is shutting down...")
		srv.Shutdown(context.Background())
	}()

	log.Print("listening...")
	srv.ListenAndServe()
	log.Print("Done")
}
