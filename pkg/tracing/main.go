package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/informatik-mannheim/cmg-ss2025/pkg/tracing/tracing"
)

// Just an example main function to test the tracing setup
func main() {
	jaeger := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if jaeger == "" {
		log.Fatal("Environment variable OTEL_EXPORTER_OTLP_ENDPOINT is not set")
	}

	// Replace with your actual Jaeger endpoint!
	shutdown, err := tracing.Init("test-service", jaeger)
	if err != nil {
		log.Fatal("Tracing init failed:", err)
	}
	defer shutdown(context.Background())

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello from Jaeger-traced service!")
	})

	log.Println("Listening on :8080...")
	http.ListenAndServe(":8080", tracing.Middleware(mux))
}
