package core

import (
	"log"

	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/model"
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/ports"
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/utils"
)

type JobSchedulerService struct {
	Notifier ports.Notifier
}

var _ ports.JobScheduler = (*JobSchedulerService)(nil)

func NewJobSchedulerService(notifier ports.Notifier) *JobSchedulerService {
	return &JobSchedulerService{
		Notifier: notifier,
	}
}

func (js *JobSchedulerService) ScheduleJob() error {
	log.Printf("Scheduling job...\n")

	// 1. Get Jobs
	jobs, err := js.Notifier.GetJobs()
	if err != nil {
		log.Printf("Error while getting jobs, aborting job-schedule: %v\n", err)
		return err
	}
	log.Printf("Jobs: %v\n", jobs)

	// 2. Get Workers
	workers, err := js.Notifier.GetWorkers()
	if err != nil {
		log.Printf("Error while getting workers, aborting job-schedule: %v\n", err)
		return err
	}
	log.Printf("Workers: %v\n", workers)

	// 3. Check already assigned
	jobs, workers, err = js.checkAlreadyAssignedJobs(jobs, workers)
	if err != nil {
		log.Printf("Error while checking already assigned jobs, aborting job-schedule: %v\n", err)
		return err
	}

	// 4. Get Carbon Intensity Data
	carbonIntensities, err := js.getCarbonIntensities(workers)
	if err != nil {
		log.Printf("Error while getting carbon intensities, aborting job-schedule: %v\n", err)
		return err
	}

	// 5. Assign Jobs to Workers
	err = js.assignJobs(jobs, workers, carbonIntensities)
	if err != nil {
		log.Printf("Error while assigning jobs, aborting job-schedule: %v\n", err)
		return err
	}

	log.Printf("Jobs assigned successfully.\n")

	return nil
}

// https://planka.123.123.123.123:8080

func (js *JobSchedulerService) checkAlreadyAssignedJobs(jobs []model.Job, workers []model.Worker) ([]model.Job, []model.Worker, error) {
	var unassignedJobs []model.Job
	for _, job := range jobs {
		if job.Status == model.JobStatusScheduled && checkAlreadyAssignedWorker(workers, job) {
			err := js.Notifier.AssignWorker(ports.UpdateWorker{
				ID: job.WorkerID,
			})
			if err != nil {
				return nil, nil, err
			}
		} else {
			unassignedJobs = append(unassignedJobs, job)
		}
	}
	unassignedWorkers := utils.Filter(workers, func(worker model.Worker) bool {
		return worker.Status == model.WorkerStatusAvailable
	})

	return unassignedJobs, unassignedWorkers, nil
}

func checkAlreadyAssignedWorker(workers []model.Worker, job model.Job) bool {
	checkWorker := func(worker model.Worker) bool {
		return worker.Id == job.WorkerID
	}
	return utils.Some(workers, checkWorker)
}

func (js *JobSchedulerService) getCarbonIntensities(workers []model.Worker) ([]model.CarbonIntensityResponse, error) {
	zones := utils.Map(workers, func(worker model.Worker) string {
		return worker.Zone
	})
	return js.Notifier.GetCarbonIntensities(zones)
}

func (js *JobSchedulerService) assignJobs(jobs []model.Job, workers []model.Worker, carbons []model.CarbonIntensityResponse) error {
	// May execute some complex algorithm to assign jobs, but for now we just assign the first available worker to the job

	var jobIndex int = 0
	for _, worker := range workers {
		if worker.Status != model.WorkerStatusAvailable {
			continue
		}

		for jobIndex < len(jobs) && jobs[jobIndex].Status != model.JobStatusScheduled {
			jobIndex++
		}
		if jobIndex >= len(jobs) {
			break
		}

		var jobPayload ports.UpdateJob = ports.UpdateJob{
			ID:              jobs[jobIndex].ID,
			WorkerID:        worker.Id,
			ComputeZone:     worker.Zone,
			CarbonIntensity: 0.0, // FIXME: mock-value for now
			CarbonSavings:   0.0, // FIXME: mock-value for now
		}
		err := js.Notifier.AssignJob(jobPayload)
		if err != nil {
			return err
		}

		// Increment so it does not get assigned again
		jobIndex++

		var workerPayload ports.UpdateWorker = ports.UpdateWorker{
			ID: worker.Id,
		}
		err = js.Notifier.AssignWorker(workerPayload)
		if err != nil {
			return err
		}
	}

	return nil
}
