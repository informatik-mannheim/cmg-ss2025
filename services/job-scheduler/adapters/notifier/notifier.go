package notifier

import (
	"log"
	"time"

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

func (n *HttpNotifier) GetJobs() ([]model.Job, error) {
	params := map[string]string{
		"param1": "value1",
		"param2": "value2",
	}
	containerImage := model.ContainerImage{
		Name:    "image1",
		Version: "1.0",
	}
	jobs := []model.Job{
		{
			ID:                   "job1",
			UserID:               "1234",
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
			JobName:              "Job 1",
			Image:                containerImage,
			AdjustmentParameters: params,
			CreationZone:         "DE",
			WorkerID:             "",
			ComputeZone:          "",
			CarbonIntensity:      -1,
			CarbonSaving:         -1,
			Result:               "",
			ErrorMessage:         "",
			Status:               model.JobStatusQueued,
		}, {
			ID:                   "job2",
			UserID:               "1234",
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
			JobName:              "Job 2",
			Image:                containerImage,
			AdjustmentParameters: params,
			CreationZone:         "DE",
			WorkerID:             "",
			ComputeZone:          "",
			CarbonIntensity:      -1,
			CarbonSaving:         -1,
			Result:               "",
			ErrorMessage:         "",
			Status:               model.JobStatusQueued,
		},
	}
	return jobs, nil
}

func (n *HttpNotifier) GetWorkers() ([]model.Worker, error) {
	workers := []model.Worker{
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

func (n *HttpNotifier) GetCarbonIntensities(zones []string) ([]model.CarbonIntensityResponse, error) {
	carbonIntensities := []model.CarbonIntensityResponse{
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
