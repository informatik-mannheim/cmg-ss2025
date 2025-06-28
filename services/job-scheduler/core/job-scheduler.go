package core

import (
	"github.com/google/uuid"
	"github.com/informatik-mannheim/cmg-ss2025/pkg/logging"
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
	logging.Debug("Scheduling jobs...")

	// 1. Get Jobs and workers
	// I do not know what vscode has drunk, but this value of workers is clearly used, im ignoring the warning
	jobs, workers, err := js.getJobsAndWorkers()
	if err != nil {
		return err
	}
	if jobs == nil || workers == nil {
		return nil // Nothing to schedule, just abort
	}

	// 2. Reassign Workers if any
	assignedJobs := GetAlreadyAssigned(jobs, workers)
	unassignedJobs := js.reassignWorkers(assignedJobs)
	jobs, workers = GetAllUnassigned(jobs, unassignedJobs, workers)

	// 3. Get Carbon Intensity Data
	zones := GetCarbonZones(jobs, workers)
	carbons, err := js.getCarbonIntensities(zones)
	if err != nil {
		return err
	}
	if carbons == nil {
		return nil // Nothing to schedule, just abort
	}

	// 4. Distribute Jobs
	jobUpdates := DistributeJobs(jobs, workers, carbons)

	// 5. Assign Jobs
	err = js.assignJobsToWorkers(jobUpdates)
	if err != nil {
		return err
	}

	logging.Debug("Job scheduling completed successfully.")
	return nil
}

func (js *JobSchedulerService) getJobsAndWorkers() ([]ports.Job, []ports.Worker, error) {
	jobs, err := js.JobAdapter.GetJobs()
	if err != nil {
		logging.Error("Error getting jobs: %v", err)
		return nil, nil, err
	}
	if len(jobs) == 0 {
		logging.Debug("No jobs available, nothing to schedule.")
		return nil, nil, nil
	}

	workers, err := js.WorkerAdapter.GetWorkers()
	if err != nil {
		logging.Error("Error getting workers: %v", err)
		return nil, nil, err
	}
	if len(workers) == 0 {
		logging.Debug("No workers available, cannot schedule jobs.")
		return nil, nil, nil
	}

	return jobs, workers, nil
}

func (js *JobSchedulerService) getCarbonIntensities(zones []string) (ports.CarbonIntensityResponse, error) {
	carbons, err := js.CarbonIntensityAdapter.GetCarbonIntensities(zones)
	if err != nil {
		logging.Error("Error getting carbon intensity data: %v", err)
		return nil, err
	}
	if len(carbons) == 0 {
		logging.Debug("No carbon intensity data available for the specified zones: %v", zones)
		return nil, nil
	}
	return carbons, nil
}

func (js *JobSchedulerService) assignJobsToWorkers(jobs []ports.UpdateJob) error {
	for _, job := range jobs {
		err := js.JobAdapter.AssignJob(job)
		if err != nil {
			logging.Error("Error assigning job %s to worker %s: %v", job.ID, job.WorkerID, err)
			return err
		}

		workerUpdate := ports.UpdateWorker{
			ID: job.WorkerID,
		}
		err = js.WorkerAdapter.AssignWorker(workerUpdate)
		if err != nil {
			logging.Error("Error updating worker %s for job %s: %v", job.WorkerID, job.ID, err)
			return err
		}
	}
	return nil
}

// returns all jobs that could not be assigned to a worker for whatever reason, those are then considered
// "not assigned" and go back into the pool
func (js *JobSchedulerService) reassignWorkers(jobs []ports.Job) []ports.Job {
	var unassignedJobs []ports.Job

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

	return unassignedJobs
}
