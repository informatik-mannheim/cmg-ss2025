package worker

import (
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

func (d *Daemon) StartHeartbeatLoop() {
	if err := d.client.Register(d.cfg.WorkerID, d.cfg.Key, d.cfg.Location); err != nil {
		fmt.Println("Registration failed:", err)
		return
	}
	fmt.Println("Worker registered successfully.")

	ticker := time.NewTicker(time.Duration(d.cfg.HeartbeatIntervalSeconds) * time.Second)
	defer ticker.Stop()

	for {
		if err := d.client.SendHeartbeat(d.cfg.WorkerID, "AVAILABLE"); err != nil {
			fmt.Println("Heartbeat failed:", err)
		} else {
			fmt.Println("Heartbeat sent.")
		}
		<-ticker.C
	}
}
