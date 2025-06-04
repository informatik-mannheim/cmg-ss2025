package core_test

import (
	"testing"

	carbonintensity "github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/adapters/carbon-intensity"
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/adapters/job"
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/adapters/worker"
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/core"
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/ports"
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/utils"
)

// No table driven tests since the Mock Data should cover all cases
func TestGetAlreadyAssigned(t *testing.T) {
	mockJobs := job.MockJobs
	mockWorkers := worker.MockWorkers

	result := core.GetAlreadyAssigned(nil, nil)
	if len(result) != 0 {
		t.Errorf("Expected 0 already assigned workers, got %d", len(result))
	}
	result = core.GetAlreadyAssigned([]ports.Job{}, []ports.Worker{})
	if len(result) != 0 {
		t.Errorf("Expected 0 already assigned workers, got %d", len(result))
	}

	alreadyAssigned := core.GetAlreadyAssigned(mockJobs, mockWorkers)

	expectedAssigned := []ports.Job{
		mockJobs[4],
	}

	if len(alreadyAssigned) != len(expectedAssigned) {
		t.Errorf("Expected %d already assigned workers, got %d", len(expectedAssigned), len(alreadyAssigned))
	}

	for i, job := range alreadyAssigned {
		if job.ID != expectedAssigned[i].ID {
			t.Errorf("Expected job ID %s, got %s", expectedAssigned[i].ID, job.ID)
		}
	}
}

// No table driven tests since the Mock Data should cover all cases
func TestGetAllUnassigned(t *testing.T) {
	mockJobs := job.MockJobs
	mockWorkers := worker.MockWorkers

	resultJobs, resultWorkers := core.GetAllUnassigned(nil, nil, nil)
	if len(resultJobs) != 0 {
		t.Errorf("Expected 0 unassigned jobs, got %d", len(resultJobs))
	}
	if len(resultWorkers) != 0 {
		t.Errorf("Expected 0 unassigned workers, got %d", len(resultWorkers))
	}

	resultJobs, resultWorkers = core.GetAllUnassigned([]ports.Job{}, []ports.Job{}, []ports.Worker{})
	if len(resultJobs) != 0 {
		t.Errorf("Expected 0 unassigned jobs, got %d", len(resultJobs))
	}
	if len(resultWorkers) != 0 {
		t.Errorf("Expected 0 unassigned workers, got %d", len(resultWorkers))
	}

	notAssigned := []ports.Job{
		mockJobs[0],
	}
	unassignedJobs, unassignedWorkers := core.GetAllUnassigned(mockJobs, notAssigned, mockWorkers)

	if len(unassignedJobs) != 4 {
		t.Errorf("Expected 4 unassigned jobs, got %d", len(unassignedJobs))
	}
	if len(unassignedWorkers) != 4 {
		t.Errorf("Expected 4 unassigned workers, got %d", len(unassignedWorkers))
	}

	expectedUnassignedJobs := []ports.Job{
		mockJobs[0],
		mockJobs[1],
		mockJobs[2],
		mockJobs[3],
	}
	expectedUnassignedWorkers := []ports.Worker{
		mockWorkers[0],
		mockWorkers[2],
		mockWorkers[3],
		mockWorkers[4],
	}

	for i, job := range unassignedJobs {
		if job.ID != expectedUnassignedJobs[i].ID {
			t.Errorf("Expected job ID %s, got %s", expectedUnassignedJobs[i].ID, job.ID)
		}
	}
	for i, worker := range unassignedWorkers {
		if worker.Id != expectedUnassignedWorkers[i].Id {
			t.Errorf("Expected worker ID %s, got %s", expectedUnassignedWorkers[i].Id, worker.Id)
		}
	}
}

// No table driven tests since the Mock Data should cover all cases
func TestGetCarbonZones(t *testing.T) {
	mockJobs := job.MockJobs
	mockWorkers := worker.MockWorkers

	result := core.GetCarbonZones(nil, nil)
	if len(result) != 0 {
		t.Errorf("Expected 0 carbon zones, got %d", len(result))
	}
	result = core.GetCarbonZones([]ports.Job{}, []ports.Worker{})
	if len(result) != 0 {
		t.Errorf("Expected 0 carbon zones, got %d", len(result))
	}

	result = core.GetCarbonZones(mockJobs, mockWorkers)
	expectedZones := make(map[string]struct{})
	expectedZones["JP"] = struct{}{}
	expectedZones["CH"] = struct{}{}
	expectedZones["FR"] = struct{}{}
	expectedZones["DE"] = struct{}{}
	expectedZones["US"] = struct{}{}

	if len(result) != len(expectedZones) {
		t.Errorf("Expected %d carbon zones, got %d", len(expectedZones), len(result))
	}

	for _, zone := range result {
		if _, ok := expectedZones[zone]; !ok {
			t.Errorf("Expected zones does not include %s", zone)
		}
	}
}

// No table driven tests since the Mock Data should cover all cases
func TestSortCarbonData(t *testing.T) {
	mockZones := carbonintensity.MockCarbons
	mockZones = append(mockZones, mockZones[0])

	result := core.SortCabonData(nil)
	if len(result) != 0 {
		t.Errorf("Expected 0 carbon data, got %d", len(result))
	}
	result = core.SortCabonData([]ports.CarbonIntensityData{})
	if len(result) != 0 {
		t.Errorf("Expected 0 carbon data, got %d", len(result))
	}

	result = core.SortCabonData(mockZones)

	expectedZones := []ports.CarbonIntensityData{
		mockZones[4],
		mockZones[1],
		mockZones[3],
		mockZones[2],
		mockZones[0],
		mockZones[0],
	}

	if len(result) != len(expectedZones) {
		t.Errorf("Expected %d carbon data, got %d", len(expectedZones), len(result))
	}

	for i, zone := range result {
		if zone != expectedZones[i] {
			t.Errorf("Expected zone %s, got %s on position %d", expectedZones[i].Zone, zone.Zone, i)
		}
	}

}

// No table driven tests since the Mock Data should cover all cases
func TestPrepareDistributionData(t *testing.T) {
	mockJobs := job.MockJobs
	mockWorkers := worker.MockWorkers
	mockCarbons := carbonintensity.MockCarbons

	resultJobs, resultWorkers, resultCarbons := core.PrepareDistributionData(nil, nil, nil)
	if len(resultJobs) != 0 {
		t.Errorf("Expected 0 jobs, got %d", len(resultJobs))
	}
	if len(resultWorkers) != 0 {
		t.Errorf("Expected 0 workers, got %d", len(resultWorkers))
	}
	if len(resultCarbons) != 0 {
		t.Errorf("Expected 0 carbons, got %d", len(resultCarbons))
	}

	resultJobs, resultWorkers, resultCarbons = core.PrepareDistributionData([]ports.Job{}, []ports.Worker{}, []ports.CarbonIntensityData{})
	if len(resultJobs) != 0 {
		t.Errorf("Expected 0 jobs, got %d", len(resultJobs))
	}
	if len(resultWorkers) != 0 {
		t.Errorf("Expected 0 workers, got %d", len(resultWorkers))
	}
	if len(resultCarbons) != 0 {
		t.Errorf("Expected 0 carbons, got %d", len(resultCarbons))
	}

	sortedCarbons := core.SortCabonData(mockCarbons)

	resultJobs, resultWorkers, resultCarbons = core.PrepareDistributionData(mockJobs, mockWorkers, sortedCarbons)

	if len(resultJobs) != 5 {
		t.Errorf("Expected 5 jobs, got %d", len(resultJobs))
	}
	if len(resultWorkers) != 5 {
		t.Errorf("Expected 5 workers, got %d", len(resultWorkers))
	}
	if len(resultCarbons) != 5 {
		t.Errorf("Expected 5 carbons, got %d", len(resultCarbons))
	}

	expectedJobs := []ports.Job{
		mockJobs[1],
		mockJobs[4],
		mockJobs[2],
		mockJobs[0],
		mockJobs[3],
	}
	expectedWorkers := []ports.Worker{
		mockWorkers[1],
		mockWorkers[4],
		mockWorkers[2],
		mockWorkers[0],
		mockWorkers[3],
	}
	expectedCarbons := make(map[string]float64)
	expectedCarbons["JP"] = 50
	expectedCarbons["CH"] = 5
	expectedCarbons["FR"] = 20
	expectedCarbons["DE"] = 100
	expectedCarbons["US"] = 10

	if len(resultCarbons) != len(expectedCarbons) {
		t.Errorf("Expected %d carbons, got %d", len(expectedCarbons), len(resultCarbons))
	}
	for i, job := range resultJobs {
		if job.ID != expectedJobs[i].ID {
			t.Errorf("Expected job ID %s, got %s", expectedJobs[i].ID, job.ID)
		}
	}
	for i, worker := range resultWorkers {
		if worker.Id != expectedWorkers[i].Id {
			t.Errorf("Expected worker ID %s, got %s", expectedWorkers[i].Id, worker.Id)
		}
	}
	for carbon, value := range resultCarbons {
		if value != expectedCarbons[carbon] {
			t.Errorf("Expected carbon %s to be %f, got %f", carbon, expectedCarbons[carbon], value)
		}
	}

}

// No table driven tests since the Mock Data should cover all cases
func TestDistributeJobs(t *testing.T) {
	mockJobs := job.MockJobs
	mockWorkers := worker.MockWorkers
	mockCarbons := carbonintensity.MockCarbons

	result := core.DistributeJobs(nil, nil, nil)
	if len(result) != 0 {
		t.Errorf("Expected 0 job updates, got %d", len(result))
	}
	result = core.DistributeJobs([]ports.Job{}, []ports.Worker{}, []ports.CarbonIntensityData{})
	if len(result) != 0 {
		t.Errorf("Expected 0 job updates, got %d", len(result))
	}

	unassignedJob := []ports.Job{mockJobs[0]}
	unassignedJobs, unassignedWorkers := core.GetAllUnassigned(mockJobs, unassignedJob, mockWorkers)

	result = core.DistributeJobs(unassignedJobs, unassignedWorkers, mockCarbons)
	expectedJobUpdates := []ports.UpdateJob{
		{
			ID:              utils.Uuid4,
			WorkerID:        utils.Uuid1,
			ComputeZone:     "JP",
			CarbonIntensity: 50,
			CarbonSavings:   50,
		},
		{
			ID:              utils.Uuid1,
			WorkerID:        utils.Uuid3,
			ComputeZone:     "FR",
			CarbonIntensity: 20,
			CarbonSavings:   80,
		},
		{
			ID:              utils.Uuid3,
			WorkerID:        utils.Uuid5,
			ComputeZone:     "US",
			CarbonIntensity: 10,
			CarbonSavings:   40,
		},
	}

	if len(result) != len(expectedJobUpdates) {
		t.Errorf("Expected %d job updates, got %d", len(expectedJobUpdates), len(result))
	}

	for i, jobUpdate := range result {
		if jobUpdate != expectedJobUpdates[i] {
			t.Errorf("Expected job update %v, got %v", expectedJobUpdates[i], jobUpdate)
		}
	}
}
