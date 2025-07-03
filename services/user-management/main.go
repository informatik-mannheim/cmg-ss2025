package main

import (
	"context"
	"github.com/informatik-mannheim/cmg-ss2025/pkg/logging"
	"github.com/informatik-mannheim/cmg-ss2025/pkg/tracing/tracing"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	auth0adapter "github.com/informatik-mannheim/cmg-ss2025/services/user-management/adapters/auth"
	handler "github.com/informatik-mannheim/cmg-ss2025/services/user-management/adapters/handler-http"
	"github.com/informatik-mannheim/cmg-ss2025/services/user-management/core"
)

func main() {
	logging.Init("user-management")
	logging.Debug("Starting User Management Service")

	jaeger := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if jaeger == "" {
		logging.Error("Environment variable OTEL_EXPORTER_OTLP_ENDPOINT is not set")
	}

	shutdown, err := tracing.Init("user-management", jaeger)
	if err != nil {
		logging.Error("Tracing init failed:", err)
	}
	defer shutdown(context.Background())

	ctx := context.Background()
	useLive := os.Getenv("USE_LIVE") == "true"

	authAdapter := auth0adapter.New(useLive)
	authService := core.NewAuthService(authAdapter)

	isAdminFn := func(secret string) bool {
		expected := os.Getenv("ADMIN_SECRET_HASH")
		return core.IsAdminSecret(secret, expected)
	}

	httpHandler := handler.New(authService, useLive, isAdminFn)

	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/auth/register", httpHandler.RegisterHandler).Methods("POST")
	muxRouter.HandleFunc("/auth/login", httpHandler.LoginHandler).Methods("POST")

	// Wrap router with tracing middleware
	tracingHandler := tracing.Middleware(muxRouter)

	server := &http.Server{
		Addr:    ":8080",
		Handler: tracingHandler,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logging.Error("Server error: ", err)
			os.Exit(1)
		}
	}()

	<-stop

	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logging.Error("Shutdown failed: ", err)
	} else {
		logging.Debug("Server shut down gracefully")
	}
}
