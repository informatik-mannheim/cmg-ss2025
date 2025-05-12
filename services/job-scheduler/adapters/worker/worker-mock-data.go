package worker

import (
	"github.com/google/uuid"
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/model"
)

var MockWorkers = []model.Worker{
	{
		Id:     uuid.New(),
		Status: model.WorkerStatusAvailable,
		Zone:   "FR",
	},
	{
		Id:     uuid.New(),
		Status: model.WorkerStatusAvailable,
		Zone:   "CH",
	},
}
