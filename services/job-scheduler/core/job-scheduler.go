package core

import (
	"log"

	"github.com/google/uuid"
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/model"
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/ports"
)

type JobSchedulerService struct {
	JobAdapter             ports.JobAdapter
	WorkerAdapter          ports.WorkerAdapter
	CarbonIntensityAdapter ports.CarbonIntensityAdapter
	Notifier               ports.Notifier
}

var _ ports.JobScheduler = (*JobSchedulerService)(nil)

func NewJobSchedulerService(notifier ports.Notifier) *JobSchedulerService {
	return &JobSchedulerService{
		Notifier: notifier,
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
	assignedJobs := getAlreadyAssigned(jobs)
	unassignedJobs := js.reassignWorkers(assignedJobs)
	jobs, workers = getAllUnassigned(jobs, unassignedJobs, workers)

	// 3. Get Carbon Intensity Data
	zones := getCarbonZones(jobs, workers)
	carbons, err := js.getCarbonIntensities(zones)
	if err != nil {
		return err
	}

	// 4. Distribute Jobs
	jobUpdates := distributeJobs(jobs, workers, carbons)

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

	workers, err := js.WorkerAdapter.GetWorkers()
	if err != nil {
		log.Printf("Error getting workers: %v\n", err)
		return nil, nil, err
	}

	return jobs, workers, nil
}

func (js *JobSchedulerService) getCarbonIntensities(zones []string) (model.CarbonIntensityResponse, error) {
	carbons, err := js.CarbonIntensityAdapter.GetCarbonIntensities(zones)
	if err != nil {
		log.Printf("Error getting carbon intensity data: %v\n", err)
		return nil, err
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
			err = js.Notifier.NotifyWorkerAssignmentFailed(job.ID, job.WorkerID)
			if err != nil {
				log.Printf("Error notifying worker assignment failed: %v\n", err)
				return err
			}
			return err
		}
		err = js.Notifier.NotifyAssignment(job.ID, job.WorkerID)
		if err != nil {
			log.Printf("Error notifying assignment: %v\n", err)
			return err
		}
	}
	return nil
}

// returns all jobs that could not be assigned to a worker for whatever reason, those are then considered
// "not assigned" and go back into the pool
func (js *JobSchedulerService) reassignWorkers(jobs []model.Job) []model.Job {
	var unassignedJobs []model.Job

	for _, job := range jobs {
		uuid, err := uuid.Parse(job.WorkerID)
		if err != nil {
			unassignedJobs = append(unassignedJobs, job)
			continue
		}

		if err := js.WorkerAdapter.AssignWorker(ports.UpdateWorker{ID: uuid}); err != nil {
			unassignedJobs = append(unassignedJobs, job)
			continue
		}

		if err := js.Notifier.NotifyAssigmentCorrection(job.ID, uuid); err != nil {
			log.Printf("Error notifying assignment correction: %v\n", err)
		}
	}

	return unassignedJobs
}
