package client_http

import (
	"fmt"
	"github.com/informatik-mannheim/cmg-ss2025/pkg/logging"
	"net/http"
	"os"
)

func PingJobScheduler() {
	url := fmt.Sprintf("%s/ping", os.Getenv(JOB_SCHEDULER_SERVICE))
	resp, err := http.Get(url)
	if err != nil {
		logging.Warn("Error pinging job-scheduler:", err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		logging.Warn("Error pinging job-scheduler:", resp.StatusCode)
	}
	defer resp.Body.Close()
}
