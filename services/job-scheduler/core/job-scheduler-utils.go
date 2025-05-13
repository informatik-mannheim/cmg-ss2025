package core

import (
	"slices"

	"github.com/google/uuid"
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/model"
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/ports"
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/utils"
)

func getAlreadyAssigned(jobs []model.Job) []model.Job {
	return utils.Filter(jobs, func(job model.Job) bool {
		return job.Status == model.JobStatusScheduled
	})
}

func getAllUnassigned(jobs, unassignedJobs []model.Job, workers []model.Worker) ([]model.Job, []model.Worker) {
	unassignedJobsMap := make(map[model.Job]struct{})
	for _, job := range unassignedJobs {
		unassignedJobsMap[job] = struct{}{}
	}

	assignedWorkersMap := make(map[uuid.UUID]struct{})

	jobResult := utils.Filter(jobs, func(job model.Job) bool {
		_, exists := unassignedJobsMap[job]
		isJobUnassigned := job.Status == model.JobStatusQueued || (job.Status == model.JobStatusScheduled && exists)

		if !isJobUnassigned && job.Status == model.JobStatusScheduled {
			uuid, err := uuid.Parse(job.WorkerID)
			if err == nil {
				assignedWorkersMap[uuid] = struct{}{}
			}
		}

		return isJobUnassigned
	})

	workerResult := utils.Filter(workers, func(worker model.Worker) bool {
		_, exists := assignedWorkersMap[worker.Id]
		return !exists
	})

	return jobResult, workerResult
}

func getCarbonZones(unassignedJobs []model.Job, unassignedWorkers []model.Worker) []string {
	zones := make(map[string]struct{})
	for _, job := range unassignedJobs {
		zones[job.ComputeZone] = struct{}{}
	}
	for _, worker := range unassignedWorkers {
		zones[worker.Zone] = struct{}{}
	}

	zonesList := make([]string, 0, len(zones))
	for zone := range zones {
		if zone == "" {
			continue
		}
		zonesList = append(zonesList, zone)
	}

	return zonesList
}

// meant are: jobs = unassigned jobs, workers = unassigned workers
func distributeJobs(jobs []model.Job, workers []model.Worker, carbons []model.CarbonIntensityData) []ports.UpdateJob {
	// small -> big
	sortedCarbons := sortCabonData(carbons)

	sortedJobs, sortedWorkers, carbonsMap := prepareDistributionData(jobs, workers, sortedCarbons)

	jobsIndex := len(sortedJobs) - 1
	workersIndex := len(sortedWorkers) - 1
	jobUpdates := make([]ports.UpdateJob, 0)

	for jobsIndex >= 0 && workersIndex >= 0 {
		job := sortedJobs[jobsIndex]
		worker := sortedWorkers[workersIndex]

		if job.CreationZone == worker.Zone {
			jobsIndex--
			continue
		}

		jobUpdate := ports.UpdateJob{
			ID:              job.ID,
			WorkerID:        worker.Id,
			ComputeZone:     worker.Zone,
			CarbonIntensity: carbonsMap[worker.Zone],
			CarbonSavings:   carbonsMap[job.CreationZone] - carbonsMap[worker.Zone],
		}
		jobUpdates = append(jobUpdates, jobUpdate)

		jobsIndex--
		workersIndex--
	}

	return jobUpdates
}

func sortCabonData(carbons []model.CarbonIntensityData) []model.CarbonIntensityData {
	slices.SortFunc(carbons, func(i, j model.CarbonIntensityData) int {
		if i.CarbonIntensity < j.CarbonIntensity {
			return -1
		} else if i.CarbonIntensity > j.CarbonIntensity {
			return 1
		}
		return 0
	})
	return carbons
}

func prepareDistributionData(jobs []model.Job, workers []model.Worker, sortedCarbons []model.CarbonIntensityData) ([]model.Job, []model.Worker, map[string]float64) {
	carbonsMap := make(map[string]float64)
	sortedJobs := make([]model.Job, 0, len(jobs))
	sortedWorkers := make([]model.Worker, 0, len(workers))

	for _, carbon := range sortedCarbons {
		carbonsMap[carbon.Zone] = carbon.CarbonIntensity
		zoneJobs := utils.Filter(jobs, func(job model.Job) bool {
			return job.CreationZone == carbon.Zone
		})
		sortedJobs = append(sortedJobs, zoneJobs...)
		zoneWorkers := utils.Filter(workers, func(worker model.Worker) bool {
			return worker.Zone == carbon.Zone
		})
		sortedWorkers = append(sortedWorkers, zoneWorkers...)
	}
	return sortedJobs, sortedWorkers, carbonsMap
}
