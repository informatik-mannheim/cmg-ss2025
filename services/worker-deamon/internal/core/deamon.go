package core

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"worker-daemon/internal/config"

	"worker-daemon/internal/ports"
)

type Daemon struct {
	cfg          config.Config
	api          ports.WorkerGateway
	workerID     string
	token        string
	currentJobID string
}

func NewDaemon(cfg config.Config, api ports.WorkerGateway) *Daemon {
	return &Daemon{cfg: cfg, api: api}
}

func (d *Daemon) StartHeartbeatLoop(ctx context.Context) {
	w, err := d.api.Register(d.cfg.Key, d.cfg.Zone)
	if err != nil {
		fmt.Println("Registration failed:", err)
		return
	}
	d.workerID = w.ID
	d.token = w.Token

	fmt.Println("Worker registered successfully.", w)

	ticker := time.NewTicker(time.Duration(d.cfg.HeartbeatIntervalSeconds) * time.Second)
	defer ticker.Stop()

	var processing int32 // 0 = not processing, 1 = processing

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Heartbeat loop stopped.")
			return

		case <-ticker.C:
			var status string
			if atomic.LoadInt32(&processing) == 1 {
				status = "RUNNING"
			} else {
				status = "AVAILABLE"
			}

			jobs, err := d.api.SendHeartbeat(d.workerID, status)
			if err != nil {
				fmt.Println("Heartbeat failed:", err)
				continue
			}
			fmt.Println("Heartbeat jobs:", jobs)

			// Wenn kein Job verarbeitet wird und neue Jobs reinkommen:
			if len(jobs) > 0 && atomic.CompareAndSwapInt32(&processing, 0, 1) {
				job := jobs[0]
				// Starte die Jobverarbeitung in separater Goroutine
				go func(j ports.Job) {
					d.currentJobID = job.ID
					processedJob := computeJob(j)
					err := d.api.SendResult(processedJob)
					if err != nil {
						fmt.Println("SendResult failed:", err)
					} else {
						atomic.StoreInt32(&processing, 0)
						d.currentJobID = ""
						status = "AVAILABLE"
					}

				}(job)
			} else {
				if len(jobs) < 1 {
					fmt.Println("No Jobs scheduled.", jobs, status, d.currentJobID)
				}
			}
		}
	}
}

func computeJob(job ports.Job) ports.Job {
	fmt.Printf("Processing job %s...\n", job.ID)
	time.Sleep(10 * time.Second)
	job.Status = "DONE"
	job.Result = "Result of job " + job.ID
	job.ErrorMessage = ""
	fmt.Printf("Job %s done.\n", job.ID)
	return job
}
