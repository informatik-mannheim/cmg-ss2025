package worker

import (
	"context"
	"fmt"
	"time"

	"worker-daemon/internal/config"

	"worker-daemon/internal/ports"
)

type Daemon struct {
	cfg    config.Config
	client ports.GatewayClient
}

func NewDaemon(cfg config.Config, client ports.GatewayClient) *Daemon {
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

	var processingJob *ports.Job = nil

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Heartbeat loop stopped.")
			return
		case <-ticker.C:
			var status string
			if processingJob != nil {
				status = "COMPUTING"
			} else {
				status = "AVAILABLE"
			}

			jobs, err := d.client.SendHeartbeat(d.cfg.WorkerID, status)
			if err != nil {
				fmt.Println("Heartbeat failed:", err)
				continue
			}

			// Wenn kein Job verarbeitet wird und neue Jobs reinkommen:
			if processingJob == nil && len(jobs) > 0 {
				job := jobs[0]
				processingJob = &job

				// Starte die Jobverarbeitung in separater Goroutine
				go func(j ports.Job) {
					processedJob := computeJob(j)
					err := d.client.SendResult(processedJob)
					if err != nil {
						fmt.Println("SendResult failed:", err)
					}
					// Job ist fertig, markiere als nil (bereit für nächsten Job)
					processingJob = nil
				}(job)
			}
		}
	}
}

func computeJob(job ports.Job) ports.Job {
	fmt.Printf("Processing job %s...\n", job.ID)
	time.Sleep(5 * time.Second) // Simuliere Arbeit
	job.Status = "DONE"
	job.Result = "Result of job " + job.ID
	job.ErrorMessage = ""
	fmt.Printf("Job %s done.\n", job.ID)
	return job
}
