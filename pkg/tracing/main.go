package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// Setup tracing
func InitTracer() (func(context.Context) error, error) {
	exp, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(sdktrace.WithBatcher(exp))
	otel.SetTracerProvider(tp)

	return tp.Shutdown, nil
}

// Middleware
func TracingMiddleware(next http.Handler) http.Handler {
	tracer := otel.Tracer("example-tracer")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, span := tracer.Start(r.Context(), r.URL.Path)
		defer span.End()

		// Add some attributes
		span.SetAttributes(
			attribute.String("http.method", r.Method),
			attribute.String("http.url", r.URL.String()),
		)

		// Call next with new context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Handler
func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "H-Hello traced world~!")
}

func main() {
	shutdown, err := InitTracer()
	if err != nil {
		log.Fatalf("failed to init tracer: %v", err)
	}
	defer shutdown(context.Background())

	mux := http.NewServeMux()
	mux.HandleFunc("/", helloHandler)

	// Apply tracing middleware
	tracedMux := TracingMiddleware(mux)

	log.Println("Server listening on :8080")
	http.ListenAndServe(":8080", tracedMux)
}
