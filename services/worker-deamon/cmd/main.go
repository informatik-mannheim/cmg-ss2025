package main

import (
	"log"

	"worker-daemon/internal/config"
	"worker-daemon/internal/gateway"
	"worker-daemon/internal/worker"
)

func main() {
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	client := gateway.NewClient(cfg.GatewayURL)
	daemon := worker.NewDaemon(cfg, client)

	go daemon.StartHeartbeatLoop()

	select {} // Block forever
}
