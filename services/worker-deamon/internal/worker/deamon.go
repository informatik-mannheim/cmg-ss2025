package worker

import (
	"context"
	"fmt"
	"time"

	"worker-daemon/internal/config"
	"worker-daemon/internal/gateway"
)

type Daemon struct {
	cfg    *config.Config
	client *gateway.Client
}

func NewDaemon(cfg *config.Config, client *gateway.Client) *Daemon {
	return &Daemon{cfg: cfg, client: client}
}

func (d *Daemon) StartHeartbeatLoop(ctx context.Context) {
	if err := d.client.Register(d.cfg.WorkerID, d.cfg.Key, d.cfg.Location); err != nil {
		fmt.Println("Registration failed:", err)
		return
	}
	fmt.Println("Worker registered successfully.")

	ticker := time.NewTicker(time.Duration(d.cfg.HeartbeatIntervalSeconds) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Heartbeat loop stopped.")
			return
		case <-ticker.C:
			if err := d.client.SendHeartbeat(d.cfg.WorkerID, "AVAILABLE"); err != nil {
				fmt.Println("Heartbeat failed:", err)
			} else {
				fmt.Println("Heartbeat sent.")
			}
		}
	}
}
