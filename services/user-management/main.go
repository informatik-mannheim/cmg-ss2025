package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	auth0adapter "github.com/informatik-mannheim/cmg-ss2025/services/user-management/adapters/auth"
	"github.com/informatik-mannheim/cmg-ss2025/services/user-management/adapters/handler-http"
	"github.com/informatik-mannheim/cmg-ss2025/services/user-management/adapters/notifier"
	"github.com/informatik-mannheim/cmg-ss2025/services/user-management/ports"
)

func main() {
	ctx := context.Background()
	useLive := os.Getenv("USE_LIVE") == "true"

	n := notifier.New()
	auth := auth0adapter.New(useLive, n)

	notifierFn := func() ports.Notifier {
		return n
	}

	h := handler.New(auth, useLive, handler.IsAdmin, notifierFn)

	mux := http.NewServeMux()
	mux.HandleFunc("/auth/register", h.RegisterHandler)
	mux.HandleFunc("/auth/login", h.LoginHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	n.Event("Listening on :8080", ctx)

	// Graceful shutdown handling
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
