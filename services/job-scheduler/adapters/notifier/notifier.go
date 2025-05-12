package notifier

import (
	"log"

	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/model"
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/ports"
)

type HttpNotifier struct{}

var _ ports.Notifier = (*HttpNotifier)(nil)

func NewHttpNotifier() *HttpNotifier {
	return &HttpNotifier{}
}

// FIXME: implement the notifier interface without mockdata

// --------------------------  Getters --------------------------

func (n *HttpNotifier) GetJobs() (model.GetJobsResponse, error) {
	jobs := model.GetJobsResponse{
		{
			ID:              "job1",
			CreationZone:    "DE",
			WorkerID:        "",
			ComputeZone:     "",
			CarbonIntensity: -1,
			CarbonSaving:    -1,
			Status:          model.JobStatusQueued,
		}, {
			ID:              "job2",
			CreationZone:    "DE",
			WorkerID:        "",
			ComputeZone:     "",
			CarbonIntensity: -1,
			CarbonSaving:    -1,
			Status:          model.JobStatusQueued,
		},
	}
	return jobs, nil
}

func (n *HttpNotifier) GetWorkers() (model.GetWorkersResponse, error) {
	workers := model.GetWorkersResponse{
		{
			Id:     "worker1",
			Status: model.WorkerStatusAvailable,
			Zone:   "DE",
		},
		{
			Id:     "worker2",
			Status: model.WorkerStatusAvailable,
			Zone:   "DE",
		},
	}
	return workers, nil
}

func (n *HttpNotifier) GetCarbonIntensities(zones []string) (model.CarbonIntensityResponse, error) {
	carbonIntensities := model.CarbonIntensityResponse{
		{
			Zone:            "DE",
			CarbonIntensity: 0.0,
		},
		{
			Zone:            "FR",
			CarbonIntensity: 0.0,
		},
	}
	return carbonIntensities, nil
}

// --------------------------  Setters --------------------------

func (n *HttpNotifier) AssignJob(update ports.UpdateJob) error {
	log.Printf("Assigned Job: %v\n", update)
	return nil
}

func (n *HttpNotifier) AssignWorker(update ports.UpdateWorker) error {
	log.Printf("Assigned Worker: %v\n", update)
	return nil
}
