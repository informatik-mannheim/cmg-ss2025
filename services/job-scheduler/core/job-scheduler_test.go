package core_test

import (
	"testing"

	carbonintensity "github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/adapters/carbon-intensity"
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/adapters/job"
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/adapters/worker"
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/core"
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/ports"
)

type TestJobScheduleTableRow struct {
	Description            string
	JobAdapter             ports.JobAdapter
	WorkerAdapter          ports.WorkerAdapter
	CarbonIntensityAdapter ports.CarbonIntensityAdapter
	Notifier               ports.Notifier
	ShouldError            bool
}

func TestScheduleJob(t *testing.T) {
	tests := []TestJobScheduleTableRow{
		{
			Description:            "Test with no errors",
			JobAdapter:             job.NewJobAdapterMock(false, false, false),
			WorkerAdapter:          worker.NewWorkerAdapterMock(false, false, false),
			CarbonIntensityAdapter: carbonintensity.NewCarbonIntensityAdapterMock(false, false),
			ShouldError:            false,
		},
		// -------------------------- Jobs --------------------------
		{
			Description:            "Test with get jobs error",
			JobAdapter:             job.NewJobAdapterMock(true, false, false),
			WorkerAdapter:          worker.NewWorkerAdapterMock(false, false, false),
			CarbonIntensityAdapter: carbonintensity.NewCarbonIntensityAdapterMock(false, false),
			ShouldError:            true,
		},
		{
			Description:            "Test with get jobs emtpy",
			JobAdapter:             job.NewJobAdapterMock(false, true, false),
			WorkerAdapter:          worker.NewWorkerAdapterMock(false, false, false),
			CarbonIntensityAdapter: carbonintensity.NewCarbonIntensityAdapterMock(false, false),
			ShouldError:            true,
		},
		{
			Description:            "Test with assign jobs error",
			JobAdapter:             job.NewJobAdapterMock(false, false, true),
			WorkerAdapter:          worker.NewWorkerAdapterMock(false, false, false),
			CarbonIntensityAdapter: carbonintensity.NewCarbonIntensityAdapterMock(false, false),
			ShouldError:            true,
		},
		// -------------------------- Workers --------------------------
		{
			Description:            "Test with get workers error",
			JobAdapter:             job.NewJobAdapterMock(false, false, false),
			WorkerAdapter:          worker.NewWorkerAdapterMock(true, false, false),
			CarbonIntensityAdapter: carbonintensity.NewCarbonIntensityAdapterMock(false, false),
			ShouldError:            true,
		},
		{
			Description:            "Test with get workers emtpy",
			JobAdapter:             job.NewJobAdapterMock(false, false, false),
			WorkerAdapter:          worker.NewWorkerAdapterMock(false, true, false),
			CarbonIntensityAdapter: carbonintensity.NewCarbonIntensityAdapterMock(false, false),
			ShouldError:            true,
		},
		{
			Description:            "Test with assign workers error",
			JobAdapter:             job.NewJobAdapterMock(false, false, false),
			WorkerAdapter:          worker.NewWorkerAdapterMock(false, false, true),
			CarbonIntensityAdapter: carbonintensity.NewCarbonIntensityAdapterMock(false, false),
			ShouldError:            true,
		},
		// -------------------------- Carbons --------------------------
		{
			Description:            "Test with get carbons error",
			JobAdapter:             job.NewJobAdapterMock(false, false, false),
			WorkerAdapter:          worker.NewWorkerAdapterMock(false, false, false),
			CarbonIntensityAdapter: carbonintensity.NewCarbonIntensityAdapterMock(true, false),
			ShouldError:            true,
		},
		{
			Description:            "Test with get carbons emtpy",
			JobAdapter:             job.NewJobAdapterMock(false, false, false),
			WorkerAdapter:          worker.NewWorkerAdapterMock(false, false, false),
			CarbonIntensityAdapter: carbonintensity.NewCarbonIntensityAdapterMock(false, true),
			ShouldError:            true,
		},
	}

	for _, row := range tests {
		t.Run(row.Description, func(t *testing.T) {
			jobSchedulerService := createJobSchedulerService(row)
			err := jobSchedulerService.ScheduleJob()
			if row.ShouldError && err == nil {
				t.Errorf("Expected error, got nil")
			} else if !row.ShouldError && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
		})
	}
}

func createJobSchedulerService(row TestJobScheduleTableRow) *core.JobSchedulerService {
	return core.NewJobSchedulerService(
		row.JobAdapter,
		row.WorkerAdapter,
		row.CarbonIntensityAdapter,
	)
}
