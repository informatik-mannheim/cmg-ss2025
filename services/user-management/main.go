package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	auth0adapter "github.com/informatik-mannheim/cmg-ss2025/services/user-management/adapters/auth"
	"github.com/informatik-mannheim/cmg-ss2025/services/user-management/adapters/handler-http"
	"github.com/informatik-mannheim/cmg-ss2025/services/user-management/adapters/notifier"
	"github.com/informatik-mannheim/cmg-ss2025/services/user-management/core"
	"github.com/informatik-mannheim/cmg-ss2025/services/user-management/ports"
)

func main() {
	ctx := context.Background()
	useLive := os.Getenv("USE_LIVE") == "true"

	// Create notifier
	n := notifier.New()

	// Create Auth0 adapter
	authAdapter := auth0adapter.New(useLive, n)

	// Wrap adapter in core AuthService
	authService := core.NewAuthService(authAdapter, n)

	// Admin verification function using core
	isAdminFn := func(secret string) bool {
		expected := os.Getenv("ADMIN_SECRET_HASH")
		return core.IsAdminSecret(secret, expected)
	}

	// HTTP handler with all dependencies
	handler := handler.New(authService, useLive, isAdminFn, func() ports.Notifier {
		return n
	})

	// Set up routes
	mux := mux.NewRouter()
	mux.HandleFunc("/auth/register", handler.RegisterHandler).Methods("POST")
	mux.HandleFunc("/auth/login", handler.LoginHandler).Methods("POST")

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	n.Event("Listening on :8080", ctx)

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			n.Event("Server error: "+err.Error(), ctx)
			os.Exit(1)
		}
	}()

	<-stop
	n.Event("Shutdown signal received", ctx)

	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		n.Event("Shutdown failed: "+err.Error(), ctx)
	} else {
		n.Event("Server shut down gracefully", ctx)
	}
}
