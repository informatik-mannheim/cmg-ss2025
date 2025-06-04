package interval_runner

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/ports"
)

type IntervalRunner struct {
	ctx      context.Context
	interval int
	port     string
}

var _ ports.Runner = (*IntervalRunner)(nil)

func NewIntervalRunner(ctx context.Context, interval int, port string) *IntervalRunner {
	return &IntervalRunner{
		ctx:      ctx,
		interval: interval,
		port:     port,
	}
}

func (ir *IntervalRunner) RunScheduleJob() {
	log.Printf("Job Scheduler starting with a %d second interval...\n", ir.interval)

	var duration = time.Duration(ir.interval) * time.Second

	for {
		select {
		case <-ir.ctx.Done():
			log.Println("Scheduler stopped.")
			return
		default:
			resp, err := http.Post("http://localhost:"+ir.port+"/schedule", "application/json", nil)
			if err != nil {
				log.Printf("Error scheduling job: %v\n", err)
			} else {
				resp.Body.Close()
			}
			select {
			case <-ir.ctx.Done():
				log.Println("Scheduler stopped.")
				return
			case <-time.After(duration):
				// continue to next iteration
			}
		}
	}
}
