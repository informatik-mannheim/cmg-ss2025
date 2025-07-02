package interval_runner

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/informatik-mannheim/cmg-ss2025/pkg/logging"
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/ports"
)

type IntervalRunner struct {
	ctx      context.Context
	interval int
	port     string
	secret   string
}

var _ ports.Runner = (*IntervalRunner)(nil)

func NewIntervalRunner(ctx context.Context, interval int, port string, secret string) *IntervalRunner {
	return &IntervalRunner{
		ctx:      ctx,
		interval: interval,
		port:     port,
		secret:   secret,
	}
}

func (ir *IntervalRunner) RunScheduleJob() {
	logging.Debug("Job Scheduler starting with a %d second interval...", ir.interval)

	var duration = time.Duration(ir.interval) * time.Second

	for {
		select {
		case <-ir.ctx.Done():
			logging.Debug("Received shutdown signal, stopping scheduler...")
			return
		default:
			resp, err := http.Post("http://localhost:"+ir.port+"/schedule", "text/plain", strings.NewReader(ir.secret))
			if err != nil {
				// Debug here because error logging already exists in the handler, would be redundant
				logging.Debug("Error scheduling job: %v", err)
			} else {
				resp.Body.Close()
			}
			select {
			case <-ir.ctx.Done():
				logging.Debug("Scheduler stopped.")
				return
			case <-time.After(duration):
				// continue to next iteration
			}
		}
	}
}
