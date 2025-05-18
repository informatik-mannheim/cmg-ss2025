package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	jobclient "github.com/informatik-mannheim/cmg-ss2025/services/consumer-gateway/adapters/client-http"
	handler_http "github.com/informatik-mannheim/cmg-ss2025/services/consumer-gateway/adapters/handler-http"
	"github.com/informatik-mannheim/cmg-ss2025/services/consumer-gateway/core"
)

func main() {

	jobClient := jobclient.New("http://job-service:8080")
	service := core.NewConsumerService(jobClient)

	srv := &http.Server{Addr: ":8080"}

	h := handler_http.NewHandler(service)
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
