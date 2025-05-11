package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	handler "github.com/informatik-mannheim/cmg-ss2025/services/carbon-intensity-provider/adapters/handler-http"
	notifier "github.com/informatik-mannheim/cmg-ss2025/services/carbon-intensity-provider/adapters/notifier"
	electricitymaps "github.com/informatik-mannheim/cmg-ss2025/services/carbon-intensity-provider/adapters/provider/electricity-maps"
	repo "github.com/informatik-mannheim/cmg-ss2025/services/carbon-intensity-provider/adapters/repo-in-memory"
	"github.com/informatik-mannheim/cmg-ss2025/services/carbon-intensity-provider/core"
	"github.com/informatik-mannheim/cmg-ss2025/services/carbon-intensity-provider/ports"
)

func main() {
	ctx := context.Background()

	n := notifier.New()
	r := repo.NewRepo()
	s := core.NewCarbonIntensityService(r, n)

	if os.Getenv("USE_LIVE") == "true" {
		log.Println("[Mode] Live fetch enabled")
		fetcher := electricitymaps.NewFromEnv(n)

		detailedZones, err := fetcher.AllElectricityMapZones(ctx)
		if err != nil {
			log.Fatalf("Failed to fetch zone metadata: %v", err)
		}

		// Filter to only zones we have a token for
		zones := make([]ports.Zone, 0)
		for _, z := range detailedZones {
			if _, ok := fetcher.TokenByZone[z.Code]; ok {
				zones = append(zones, ports.Zone{Code: z.Code, Name: z.Name})
			}
		}

		if err := r.StoreZones(zones, ctx); err != nil {
			n.Event("Failed to store filtered zone metadata: " + err.Error())
		}

		configuredZones, err := fetcher.GetConfiguredZones(ctx)
		if err != nil {
			log.Fatalf("Failed to get configured zones: %v", err)
		}

		go func() {
			for {
				for _, zone := range configuredZones {
					data, err := fetcher.Fetch(zone, ctx)
					if err != nil {
						n.Event("Error fetching data for zone " + zone + ": " + err.Error())
						continue
					}
					s.AddOrUpdateZone(data.Zone, data.CarbonIntensity, ctx)
				}
				time.Sleep(60 * time.Second)
			}
		}()
	} else {
		log.Println("[Mode] Using static offline test data")
		s.AddOrUpdateZone("DE", 140.5, ctx)
		s.AddOrUpdateZone("FR", 135.2, ctx)
		s.AddOrUpdateZone("US-NY-NYIS", 128.9, ctx)
	}

	httpHandler := handler.NewHandler(s, n)
	server := &http.Server{
		Addr:    ":8080",
		Handler: httpHandler,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Println("Carbon Intensity Provider running on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-stop
	log.Println("Shutdown signal received")

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctxShutdown); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server exited gracefully")
}
