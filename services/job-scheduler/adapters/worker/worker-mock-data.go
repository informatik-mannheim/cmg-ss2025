package worker

import (
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/model"
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/utils"
)

var MockWorkers = []model.Worker{
	{
		Id:     utils.Uuid1,
		Status: model.WorkerStatusAvailable,
		Zone:   "JP",
	},
	{
		Id:     utils.Uuid2,
		Status: model.WorkerStatusAvailable,
		Zone:   "CH",
	},
	{
		Id:     utils.Uuid3,
		Status: model.WorkerStatusAvailable,
		Zone:   "FR",
	},
	{
		Id:     utils.Uuid4,
		Status: model.WorkerStatusAvailable,
		Zone:   "DE",
	},
	{
		Id:     utils.Uuid5,
		Status: model.WorkerStatusAvailable,
		Zone:   "US",
	},
}
