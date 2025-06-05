package job

import (
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/ports"
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/utils"
)

var MockJobs = []ports.Job{
	{
		ID:              utils.Uuid1,
		CreationZone:    "DE",
		WorkerID:        "",
		ComputeZone:     "",
		CarbonIntensity: -1,
		CarbonSaving:    -1,
		Status:          ports.JobStatusScheduled,
	},
	{
		ID:              utils.Uuid2,
		CreationZone:    "US",
		WorkerID:        "",
		ComputeZone:     "",
		CarbonIntensity: -1,
		CarbonSaving:    -1,
		Status:          ports.JobStatusQueued,
	},
	{
		ID:              utils.Uuid3,
		CreationZone:    "JP",
		WorkerID:        "",
		ComputeZone:     "",
		CarbonIntensity: -1,
		CarbonSaving:    -1,
		Status:          ports.JobStatusQueued,
	},
	{
		ID:              utils.Uuid4,
		CreationZone:    "DE",
		WorkerID:        utils.Uuid1.String(),
		ComputeZone:     "", // Does not matter as of now
		CarbonIntensity: -1,
		CarbonSaving:    -1,
		Status:          ports.JobStatusQueued,
	},
	{
		ID:              utils.Uuid5,
		CreationZone:    "US",
		WorkerID:        utils.Uuid2.String(),
		ComputeZone:     "", // Does not matter as of now
		CarbonIntensity: -1,
		CarbonSaving:    -1,
		Status:          ports.JobStatusScheduled,
	},
}
