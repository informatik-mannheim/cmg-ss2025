package worker

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/ports"
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/utils"
)

func PutWorkerStatusEndpoint(base string, id uuid.UUID) string {
	return fmt.Sprintf("%s/workers/%s/status", base, id)
}

func GetWorkersEndpoint(base string) string {
	baseUrl := fmt.Sprintf("%s/workers", base)

	params := url.Values{}
	params.Add("status", string(ports.WorkerStatusAvailable))

	fullUrl := baseUrl + "?" + params.Encode()
	return fullUrl
}

type WorkerAdapter struct {
	baseUrl string
	client  http.Client
}

var _ ports.WorkerAdapter = (*WorkerAdapter)(nil)

func NewWorkerAdapter(client http.Client, baseUrl string) *WorkerAdapter {
	return &WorkerAdapter{
		baseUrl: baseUrl,
		client:  client,
	}
}

func (adapter *WorkerAdapter) GetWorkers() (ports.GetWorkersResponse, error) {
	// For now its kept simple and return an error as soon as it gets one, changes in Phase 3
	endpoint := GetWorkersEndpoint(adapter.baseUrl)

	// StatusCode is not relevant yet
	data, _, err := utils.GetRequest[ports.GetWorkersResponse](&adapter.client, endpoint)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (adapter *WorkerAdapter) AssignWorker(update ports.UpdateWorker) error {
	// For now its kept simple and return an error as soon as it gets one, changes in Phase 3
	endpoint := PutWorkerStatusEndpoint(adapter.baseUrl, update.ID)

	payload := ports.UpdateWorkerPayload{
		WorkerStatus: ports.WorkerStatusRunning,
	}

	// StatusCode is not relevant yet
	_, _, err := utils.PutRequest[ports.UpdateWorkerPayload, ports.UpdateWorkerResponse](&adapter.client, endpoint, payload)
	if err != nil {
		return err
	}

	return nil
}
