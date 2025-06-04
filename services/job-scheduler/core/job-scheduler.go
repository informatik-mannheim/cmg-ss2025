package core

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/model"
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/ports"
)

type JobSchedulerService struct {
	JobAdapter             ports.JobAdapter
	WorkerAdapter          ports.WorkerAdapter
	CarbonIntensityAdapter ports.CarbonIntensityAdapter
}

var _ ports.JobScheduler = (*JobSchedulerService)(nil)

func NewJobSchedulerService(
	jobAdapter ports.JobAdapter,
	workerAdapter ports.WorkerAdapter,
	carbonIntensityAdapter ports.CarbonIntensityAdapter,
) *JobSchedulerService {
	return &JobSchedulerService{
		JobAdapter:             jobAdapter,
		WorkerAdapter:          workerAdapter,
		CarbonIntensityAdapter: carbonIntensityAdapter,
	}
}

func (js *JobSchedulerService) ScheduleJob() error {
	log.Printf("Scheduling job...\n")

	// 1. Get Jobs and workers
	// I do not know what vscode has drunk, but this value of workers is clearly used, im ignoring the warning
	jobs, workers, err := js.getJobsAndWorkers()
	if err != nil {
		return err
	}

	// 2. Reassign Workers if any
	assignedJobs := GetAlreadyAssigned(jobs, workers)
	unassignedJobs, err := js.reassignWorkers(assignedJobs)
	if err != nil {
		return err
	}
	jobs, workers = GetAllUnassigned(jobs, unassignedJobs, workers)

	// 3. Get Carbon Intensity Data
	zones := GetCarbonZones(jobs, workers)
	carbons, err := js.getCarbonIntensities(zones)
	if err != nil {
		return err
	}

	// 4. Distribute Jobs
	jobUpdates := DistributeJobs(jobs, workers, carbons)

	// 5. Assign Jobs
	err = js.assignJobsToWorkers(jobUpdates)
	if err != nil {
		return err
	}

	log.Printf("Jobs assigned successfully.\n")
	return nil
}

func (js *JobSchedulerService) getJobsAndWorkers() ([]model.Job, []model.Worker, error) {
	jobs, err := js.JobAdapter.GetJobs()
	if err != nil {
		log.Printf("Error getting jobs: %v\n", err)
		return nil, nil, err
	}
	if len(jobs) == 0 {
		log.Printf("No jobs available.\n")
		return nil, nil, fmt.Errorf("no jobs available")
	}

	workers, err := js.WorkerAdapter.GetWorkers()
	if err != nil {
		log.Printf("Error getting workers: %v\n", err)
		return nil, nil, err
	}
	if len(workers) == 0 {
		log.Printf("No workers available.\n")
		return nil, nil, fmt.Errorf("no workers available")
	}

	return jobs, workers, nil
}

func (js *JobSchedulerService) getCarbonIntensities(zones []string) (model.CarbonIntensityResponse, error) {
	carbons, err := js.CarbonIntensityAdapter.GetCarbonIntensities(zones)
	if err != nil {
		log.Printf("Error getting carbon intensity data: %v\n", err)
		return nil, err
	}
	if len(carbons) == 0 {
		log.Printf("No carbon intensity data available.\n")
		return nil, fmt.Errorf("no carbon intensity data available")
	}
	return carbons, nil
}

func (js *JobSchedulerService) assignJobsToWorkers(jobs []ports.UpdateJob) error {
	for _, job := range jobs {
		err := js.JobAdapter.AssignJob(job)
		if err != nil {
			log.Printf("Error updating job: %v\n", err)
			return err
		}

		workerUpdate := ports.UpdateWorker{
			ID: job.WorkerID,
		}
		err = js.WorkerAdapter.AssignWorker(workerUpdate)
		if err != nil {
			log.Printf("Error updating worker: %v\n", err)
			return err
		}
	}
	return nil
}

// returns all jobs that could not be assigned to a worker for whatever reason, those are then considered
// "not assigned" and go back into the pool
func (js *JobSchedulerService) reassignWorkers(jobs []model.Job) ([]model.Job, error) {
	var unassignedJobs []model.Job

	for _, job := range jobs {
		// error ignored on purpose, since here can only be jobs that have an workerId
		// for an worker that was fetched, and since the Id in the worker is typesafe as
		// uuid, we can safely ignore the error since it will never happen
		uuid, _ := uuid.Parse(job.WorkerID)

		if err := js.WorkerAdapter.AssignWorker(ports.UpdateWorker{ID: uuid}); err != nil {
			unassignedJobs = append(unassignedJobs, job)
			continue
		}

	}

	return unassignedJobs, nil
}
