package main

import (
	"context"
	"github.com/informatik-mannheim/cmg-ss2025/pkg/logging"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/informatik-mannheim/cmg-ss2025/pkg/tracing/tracing"
	handler "github.com/informatik-mannheim/cmg-ss2025/services/carbon-intensity-provider/adapters/handler-http"
	notifier "github.com/informatik-mannheim/cmg-ss2025/services/carbon-intensity-provider/adapters/notifier"
	electricitymaps "github.com/informatik-mannheim/cmg-ss2025/services/carbon-intensity-provider/adapters/provider/electricity-maps"
	repo "github.com/informatik-mannheim/cmg-ss2025/services/carbon-intensity-provider/adapters/repo-in-memory"
	"github.com/informatik-mannheim/cmg-ss2025/services/carbon-intensity-provider/core"
	"github.com/informatik-mannheim/cmg-ss2025/services/carbon-intensity-provider/ports"
)

func main() {
	logging.Init("carbon-intensity-provider")
	logging.Debug("Starting Carbon Intensity Provider")

	jaeger := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if jaeger == "" {
		logging.Error("Environment variable OTEL_EXPORTER_OTLP_ENDPOINT is not set")
	}

	shutdown, err := tracing.Init("carbon-intensity-provider", jaeger)
	if err != nil {
		logging.Error("Tracing init failed:", err)
	}
	defer shutdown(context.Background())

	rootCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	n := notifier.New()
	r := repo.NewRepo(n)
	s := core.NewCarbonIntensityService(r, n)

	if os.Getenv("USE_LIVE") == "true" {
		logging.Debug("[Mode] Live fetch enabled")
		fetcher := electricitymaps.NewFromEnv(n)

		detailedZones, err := fetcher.AllElectricityMapZones(rootCtx)
		if err != nil {
			errormessage := "Failed to fetch zone metadata: " + err.Error()
			logging.Error(errormessage)
		}

		// Filter to only zones we have a token for
		zones := make([]ports.Zone, 0)
		for _, z := range detailedZones {
			if _, ok := fetcher.TokenByZone[z.Code]; ok {
				zones = append(zones, ports.Zone{Code: z.Code, Name: z.Name})
			}
		}

		if err := r.StoreZones(zones, rootCtx); err != nil {
			n.Event("Failed to store filtered zone metadata: " + err.Error())
		}

		configuredZones, err := fetcher.GetConfiguredZones(rootCtx)
		if err != nil {
			logging.Error("Failed to get configured zones: " + err.Error())
		}

		// Context-aware fetch loop
		go func(ctx context.Context) {
			ticker := time.NewTicker(60 * time.Second)
			defer ticker.Stop()

			for {
				select {
				case <-ctx.Done():
					n.Event("Fetcher loop stopped due to context cancellation")
					return
				case <-ticker.C:
					for _, zone := range configuredZones {
						data, err := fetcher.Fetch(zone, ctx)
						if err != nil {
							n.Event("Error fetching data for zone " + zone + ": " + err.Error())
							continue
						}
						s.AddOrUpdateZone(data.Zone, data.CarbonIntensity, ctx)
					}
				}
			}
		}(rootCtx)

	} else {
		logging.Debug("[Mode] Using static offline test data")
		s.AddOrUpdateZone("DE", 140.5, rootCtx)
		s.AddOrUpdateZone("FR", 135.2, rootCtx)
		s.AddOrUpdateZone("US-NY-NYIS", 128.9, rootCtx)

		offlineZones := []ports.Zone{
			{Code: "DE", Name: "Germany"},
			{Code: "FR", Name: "France"},
			{Code: "US-NY-NYIS", Name: "New York ISO"},
		}
		_ = r.StoreZones(offlineZones, rootCtx)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	server := &http.Server{Addr: ":" + port}

	httpHandler := handler.NewHandler(s, n)
	tracingHandler := tracing.Middleware(httpHandler)
	http.Handle("/", tracingHandler)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logging.Debug("Carbon Intensity Provider running on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logging.Error("Server error: ", err.Error())
		}
	}()

	<-stop
	n.Event("Shutdown signal received")
	cancel() // cancel root context

	ctxShutdown, cancelTimeout := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelTimeout()

	if err := server.Shutdown(ctxShutdown); err != nil {
		logging.Error("Server shutdown failed: " + err.Error())
	}

	n.Event("Server exited gracefully")
	logging.Debug("Server exited gracefully")
}
