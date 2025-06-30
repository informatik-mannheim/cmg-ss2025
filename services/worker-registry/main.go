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

	jwksUrl := os.Getenv("JWKS_URL")
	err := auth.InitJWKS(jwksUrl)

	if err != nil {
		logging.Error("Failed to initialize JWKS: " + err.Error())
		return
	}

	jaeger := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if jaeger == "" {
		logging.Error("Environment variable OTEL_EXPORTER_OTLP_ENDPOINT is not set")
	}

	shutdown, err := tracing.Init("worker-registry", jaeger)
	if err != nil {
		logging.Error("Tracing init failed:", err)
	}
	defer shutdown(context.Background())

	zoneClient := client.NewZoneClient(os.Getenv("CARBON_INTENSITY_PROVIDER"))

	var service *core.WorkerRegistryService
	var dbRepo ports.Repo

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	sslModeStr := os.Getenv("SSL_MODE")
	sslMode, err := strconv.ParseBool(sslModeStr)
	maxRetries := 10

	if err != nil {
		logging.Warn("Invalid SSL_MODE value, defaulting to false")
		sslMode = false
	}

	for i := 0; i < maxRetries; i++ {
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

	service = core.NewWorkerRegistryService(dbRepo, zoneClient)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	srv := &http.Server{Addr: ":" + port}

	h := handler.NewHandler(service)
	tracingHandler := tracing.Middleware(h)
	http.Handle("/", tracingHandler)

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		logging.Debug("The service is shutting down...")
		srv.Shutdown(context.Background())
	}()

	logging.Debug("Listening...")
	srv.ListenAndServe()
	logging.Debug("Done")
}
