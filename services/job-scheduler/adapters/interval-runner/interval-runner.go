package interval_runner

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
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
	logging.Debug(fmt.Sprintf("Job Scheduler starting with a %d second interval...", ir.interval))

	var duration = time.Duration(ir.interval) * time.Second

	for {
		select {
		case <-ir.ctx.Done():
			logging.Debug("Received shutdown signal, stopping scheduler...")
			return
		default:
			url := fmt.Sprintf("http://localhost:%s/schedule", ir.port)

			// error can be ignored here, since its a static request
			requestBody, _ := json.Marshal(ports.ScheduleRequest{Secret: ir.secret})
			resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))

			// Debug-Logging here because error logging already exists in the handler, would be redundant
			// (Its more fine-tuned like this)
			if err != nil {
				logging.Debug(fmt.Sprintf("Error scheduling job: %v", err))
			} else {
				var errResponse ports.ScheduleResponse

				if err := json.NewDecoder(resp.Body).Decode(&errResponse); err != nil {
					logging.Debug(fmt.Sprintf("Job-Scheduler Response: received status code '%d' with response: 'undefined'", resp.StatusCode))
				} else {
					logging.Debug(fmt.Sprintf("Job-Scheduler Response: received status code '%d' with response: '%s'", resp.StatusCode, errResponse))
				}

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
