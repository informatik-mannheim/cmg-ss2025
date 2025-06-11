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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// set port manually with "export PORT"
	job := jobclient.NewJobClient(os.Getenv("JOB_SERVICE") + ":" + port)
	user := jobclient.NewLoginClient(os.Getenv("USER_MANAGEMENT_SERVICE") + ":" + port)
	zone := jobclient.NewZoneClient(os.Getenv("CARBON_INTENSITY_PROVIDER") + ":" + port)
	service := core.NewConsumerService(job, zone, user)
	handler := handler_http.NewHandler(service)

	mux := http.NewServeMux()
	mux.Handle("/", handler)

	srv := &http.Server{Addr: ":" + port, Handler: mux}

	mux.HandleFunc("/jobs", handler.HandleCreateJobRequest)
	mux.HandleFunc("jobs/{id}/result", handler.HandleGetJobOutcomeRequest)
	mux.HandleFunc("/auth/login", handler.HandleLoginRequest)

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Print("The service is shutting down...")
		err := srv.Shutdown(context.Background())
		if err != nil {
			return
		}
	}()

	log.Print("listening on port " + port + " ...")

	err := srv.ListenAndServe()
	if err != nil {
		return
	}
	log.Print("Done")
}
