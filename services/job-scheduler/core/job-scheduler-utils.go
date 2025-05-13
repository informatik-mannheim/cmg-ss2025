package core

import (
	"slices"

	"github.com/google/uuid"
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/model"
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/ports"
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/utils"
)

func GetAlreadyAssigned(jobs []model.Job, workers []model.Worker) []model.Job {
	workerMap := make(map[uuid.UUID]struct{})
	for _, worker := range workers {
		workerMap[worker.Id] = struct{}{}
	}

	return utils.Filter(jobs, func(job model.Job) bool {
		uuid, err := uuid.Parse(job.WorkerID)
		if err != nil {
			return false
		}
		_, exists := workerMap[uuid]
		return job.Status == model.JobStatusScheduled && exists
	})
}

func GetAllUnassigned(jobs, unassignedJobs []model.Job, workers []model.Worker) ([]model.Job, []model.Worker) {
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

func GetCarbonZones(unassignedJobs []model.Job, unassignedWorkers []model.Worker) []string {
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
func DistributeJobs(jobs []model.Job, workers []model.Worker, carbons []model.CarbonIntensityData) []ports.UpdateJob {
	// small -> big
	sortedCarbons := SortCabonData(carbons)

	sortedJobs, sortedWorkers, carbonsMap := PrepareDistributionData(jobs, workers, sortedCarbons)

	jobsIndex := len(sortedJobs) - 1
	workersIndex := len(sortedWorkers) - 1
	jobUpdates := make([]ports.UpdateJob, 0)

	for jobsIndex >= 0 && workersIndex >= 0 {
		job := sortedJobs[jobsIndex]
		worker := sortedWorkers[workersIndex]

		if job.CreationZone == worker.Zone {
			workersIndex--
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

func SortCabonData(carbons []model.CarbonIntensityData) []model.CarbonIntensityData {
	copyCarbons := make([]model.CarbonIntensityData, len(carbons))
	copy(copyCarbons, carbons)
	slices.SortFunc(copyCarbons, func(i, j model.CarbonIntensityData) int {
		if i.CarbonIntensity < j.CarbonIntensity {
			return -1
		} else if i.CarbonIntensity > j.CarbonIntensity {
			return 1
		}
		return 0
	})
	return copyCarbons
}

func PrepareDistributionData(jobs []model.Job, workers []model.Worker, sortedCarbons []model.CarbonIntensityData) ([]model.Job, []model.Worker, map[string]float64) {
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
