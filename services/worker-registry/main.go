package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/informatik-mannheim/cmg-ss2025/pkg/auth"
	"github.com/informatik-mannheim/cmg-ss2025/pkg/logging"
	"github.com/informatik-mannheim/cmg-ss2025/pkg/tracing/tracing"

	client "github.com/informatik-mannheim/cmg-ss2025/services/worker-registry/adapters/clients"
	handler "github.com/informatik-mannheim/cmg-ss2025/services/worker-registry/adapters/handler-http"
	repo_pg "github.com/informatik-mannheim/cmg-ss2025/services/worker-registry/adapters/repo"
	"github.com/informatik-mannheim/cmg-ss2025/services/worker-registry/core"
	"github.com/informatik-mannheim/cmg-ss2025/services/worker-registry/ports"
)

func main() {
	logging.Init("worker-registry")
	logging.Debug("Starting Worker Registry")

	// JWKS Token
	if err := auth.InitJWKS(os.Getenv("JWKS_URL")); err != nil {
		logging.Error("Failed to initialize JWKS: " + err.Error())
		return
	}

	// Tracing
	jaeger := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	shutdown, err := tracing.Init("worker-registry", jaeger)
	if err != nil {
		logging.Error("Tracing init failed:", err)
	}
	defer shutdown(context.Background())

	// Zone-Client
	rawInterval := os.Getenv("CARBON_INTENSITY_PROVIDER_INTERVAL")
	intervalSeconds := 60

	if rawInterval != "" {
		if i, err := strconv.Atoi(rawInterval); err == nil {
			intervalSeconds = i
		}
	}

	cacheDuration := time.Duration(intervalSeconds) * time.Second

	zoneClient := client.NewZoneClient(
		os.Getenv("CARBON_INTENSITY_PROVIDER"),
		cacheDuration,
	)

	// DB-Connection
	var dbRepo ports.Repo
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	retries := 10

	sslMode, err := strconv.ParseBool(os.Getenv("SSL_MODE"))
	if err != nil {
		logging.Warn("Invalid SSL_MODE value, defaulting to false")
		sslMode = false
	}

	// Retry-Loop
	for range retries {
		r, err := repo_pg.NewRepo(dbHost, dbPort, dbUser, dbPassword, dbName, sslMode, context.Background())
		if err == nil {
			dbRepo = r
			break
		}
		logging.Warn("Postgres not ready, retrying...")
		time.Sleep(2 * time.Second)
	}
	if dbRepo == nil {
		logging.Error("Failed to connect to Postgres")
		return
	}

	service := core.NewWorkerRegistryService(dbRepo, zoneClient)
	httpHandler := handler.NewHandler(service)

	// Authn + Propagate Authorization-Header + Tracing
	mux := http.NewServeMux()
	mux.Handle("/", tracing.Middleware(auth.AuthMiddleware(PropagateAuthMiddleware(httpHandler))))

	// HTTP-Server
	srv := &http.Server{Addr: ":" + getPort(), Handler: mux}
	go shutdownOnSignal(srv)

	logging.Debug("Listening on port " + getPort())
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logging.Error("Server error", "err", err)
	}
	logging.Debug("Done")
}

func getPort() string {
	if port := os.Getenv("PORT"); port != "" {
		return port
	}
	return "8080"
}

func shutdownOnSignal(srv *http.Server) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	logging.Debug("Shutting down...")
	srv.Shutdown(context.Background())
}

// PropagateAuthMiddleware pulls the incoming Bearer-Token and adds it to the context
func PropagateAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		ctx := context.WithValue(r.Context(), "authHeader", authHeader)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
