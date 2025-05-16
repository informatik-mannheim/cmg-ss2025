package ports

import (
	"github.com/google/uuid"
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/model"
)

type UpdateWorker struct {
	ID uuid.UUID `json:"id"`
	// No status on this struct, because there is only 1 possible option,
	// so the function will set it itself.
}

type WorkerAdapter interface {
	GetWorkers() (model.GetWorkersResponse, error)
	AssignWorker(update UpdateWorker) error
}
