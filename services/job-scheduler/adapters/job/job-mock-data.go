package job

import (
	"github.com/google/uuid"
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/model"
)

var MockJobs = []model.Job{
	{
		ID:              uuid.New(),
		CreationZone:    "DE",
		WorkerID:        "",
		ComputeZone:     "",
		CarbonIntensity: -1,
		CarbonSaving:    -1,
		Status:          model.JobStatusQueued,
	},
	{
		ID:              uuid.New(),
		CreationZone:    "US",
		WorkerID:        "",
		ComputeZone:     "",
		CarbonIntensity: -1,
		CarbonSaving:    -1,
		Status:          model.JobStatusQueued,
	},
	{
		ID:              uuid.New(),
		CreationZone:    "JP",
		WorkerID:        "",
		ComputeZone:     "",
		CarbonIntensity: -1,
		CarbonSaving:    -1,
		Status:          model.JobStatusQueued,
	},
}
