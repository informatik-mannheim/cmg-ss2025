package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go daemon.StartHeartbeatLoop(ctx)

	// Graceful shutdown handling
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop // Wait for termination signal
	log.Println("Shutting down daemon...")

	cancel()

	log.Println("Shutdown complete")

	select {} // Block forever
}
