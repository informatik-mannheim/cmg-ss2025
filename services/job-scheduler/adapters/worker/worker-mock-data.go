package worker

import (
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/ports"
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/utils"
)

var MockWorkers = []ports.Worker{
	{
		Id:     utils.Uuid1,
		Status: ports.WorkerStatusAvailable,
		Zone:   "JP",
	},
	{
		Id:     utils.Uuid2,
		Status: ports.WorkerStatusAvailable,
		Zone:   "CH",
	},
	{
		Id:     utils.Uuid3,
		Status: ports.WorkerStatusAvailable,
		Zone:   "FR",
	},
	{
		Id:     utils.Uuid4,
		Status: ports.WorkerStatusAvailable,
		Zone:   "DE",
	},
	{
		Id:     utils.Uuid5,
		Status: ports.WorkerStatusAvailable,
		Zone:   "US",
	},
}
