package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	handler_http "github.com/informatik-mannheim/cmg-ss2025/services/consumer-gateway/adapters/handler-http"
	jobclient "github.com/informatik-mannheim/cmg-ss2025/services/consumer-gateway/adapters/job_client-http"
	loginclient "github.com/informatik-mannheim/cmg-ss2025/services/consumer-gateway/adapters/job_client-http"
	zoneclient "github.com/informatik-mannheim/cmg-ss2025/services/consumer-gateway/adapters/job_client-http"
	"github.com/informatik-mannheim/cmg-ss2025/services/consumer-gateway/core"
)

func main() {

	job := jobclient.New("http://job:8080")
	user := loginclient.New("http://user-management:8080")
	zone := zoneclient.New("http://carbon-intensity-provider:8080")
	service := core.NewConsumerService(job, zone, user)

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
