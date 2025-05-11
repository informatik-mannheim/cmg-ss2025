package notifier

import (
	"log"
	"time"

	CarbonIntensityProvider "github.com/informatik-mannheim/cmg-ss2025/services/carbon-intensity-provider/model"
	Job "github.com/informatik-mannheim/cmg-ss2025/services/job"
	JobScheduler "github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/ports"
	WorkerRegistry "github.com/informatik-mannheim/cmg-ss2025/services/worker-registry/ports"
)

type HttpNotifier struct{}

var _ JobScheduler.Notifier = (*HttpNotifier)(nil)

func NewHttpNotifier() *HttpNotifier {
	return &HttpNotifier{}
}

// FIXME: implement the notifier interface without mockdata

// --------------------------  Getters --------------------------

func (n *HttpNotifier) GetJobs() ([]Job.Job, error) {
	var params map[string]string = map[string]string{
		"param1": "value1",
		"param2": "value2",
	}
	var containerImage Job.ContainerImage = Job.ContainerImage{
		Name:    "image1",
		Version: "1.0",
	}
	var jobs []Job.Job = []Job.Job{
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
			Status:               Job.StatusQueued,
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
			Status:               Job.StatusQueued,
		},
	}
	return jobs, nil
}

func (n *HttpNotifier) GetWorkers() ([]WorkerRegistry.Worker, error) {
	var workers []WorkerRegistry.Worker = []WorkerRegistry.Worker{
		{
			Id:     "worker1",
			Status: WorkerRegistry.StatusAvailable,
			Zone:   "DE",
		},
		{
			Id:     "worker2",
			Status: WorkerRegistry.StatusAvailable,
			Zone:   "DE",
		},
	}
	return workers, nil
}

func (n *HttpNotifier) GetCarbonIntensities(zones []string) ([]CarbonIntensityProvider.CarbonIntensityData, error) {
	var carbonIntensities []CarbonIntensityProvider.CarbonIntensityData = []CarbonIntensityProvider.CarbonIntensityData{
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

func (n *HttpNotifier) AssignJob(update JobScheduler.UpdateJob) error {
	log.Printf("Assigned Job: %v\n", update)
	return nil
}

func (n *HttpNotifier) AssignWorker(update JobScheduler.UpdateWorker) error {
	log.Printf("Assigned Worker: %v\n", update)
	return nil
}
