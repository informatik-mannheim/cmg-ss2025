package main

import (
	"log"
	"net/http"

	"github.com/informatik-mannheim/cmg-ss2025/services/carbon-intensity-provider/api"
	"github.com/informatik-mannheim/cmg-ss2025/services/carbon-intensity-provider/core"
)

func main() {
	service := core.NewCarbonIntensityService()

	// Preload fixed zones manually as test(for Assignment II)
	service.AddOrUpdateZone("DE", 140.5)
	service.AddOrUpdateZone("FR", 135.2)
	service.AddOrUpdateZone("US-NY-NYIS", 128.9)

	// Start server
	router := api.NewHandler(service)
	log.Println("Carbon Intensity Provider running on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
