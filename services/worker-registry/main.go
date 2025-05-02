package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	handler_http "github.com/informatik-mannheim/cmg-ss2025/services/worker-registry/adapters/handler-http"
	repo "github.com/informatik-mannheim/cmg-ss2025/services/worker-registry/adapters/repo-in-memory"
	"github.com/informatik-mannheim/cmg-ss2025/services/worker-registry/core"
)

func main() {

	core := core.NewWorkerRegistryService(repo.NewRepo(), nil)

	srv := &http.Server{Addr: ":8080"}

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
