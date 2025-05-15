package core_test

import (
	"context"
	"errors"
	"testing"

	"github.com/informatik-mannheim/cmg-ss2025/services/worker-gateway/core"
	"github.com/informatik-mannheim/cmg-ss2025/services/worker-gateway/ports"
)

// --- Dummy Notifier für Tests ---
type dummyNotifier struct {
	RegisterWorkerCalled     bool
	UpdateJobCalled          bool
	UpdateWorkerStatusCalled bool
	FetchScheduledJobsCalled bool

	ReturnErr bool
}

func (d *dummyNotifier) RegisterWorker(ctx context.Context, req ports.RegisterRequest) error {
	d.RegisterWorkerCalled = true
	if d.ReturnErr {
		return errors.New("register error")
	}
	return nil
}

func (d *dummyNotifier) UpdateJob(ctx context.Context, req ports.ResultRequest) error {
	d.UpdateJobCalled = true
	if d.ReturnErr {
		return errors.New("update job error")
	}
	return nil
}

func (d *dummyNotifier) UpdateWorkerStatus(ctx context.Context, req ports.HeartbeatRequest) error {
	d.UpdateWorkerStatusCalled = true
	if d.ReturnErr {
		return errors.New("update worker status error")
	}
	return nil
}

func (d *dummyNotifier) FetchScheduledJobs(ctx context.Context) ([]ports.Job, error) {
	d.FetchScheduledJobsCalled = true
	if d.ReturnErr {
		return nil, errors.New("fetch jobs error")
	}
	return []ports.Job{
		{ID: "job1", WorkerID: "worker1", Status: "SCHEDULED"},
		{ID: "job2", WorkerID: "worker2", Status: "SCHEDULED"},
	}, nil
}

// --- Tests ---

func TestRegisterWorker_Success(t *testing.T) {
	notifier := &dummyNotifier{}
	svc := core.NewWorkerGatewayService(notifier)

	req := ports.RegisterRequest{
		ID:       "worker1",
		Key:      "secret",
		Location: "DE",
	}

	err := svc.RegisterWorker(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !notifier.RegisterWorkerCalled {
		t.Error("expected RegisterWorker to be called")
	}
}

func TestRegisterWorker_Error(t *testing.T) {
	notifier := &dummyNotifier{ReturnErr: true}
	svc := core.NewWorkerGatewayService(notifier)

	err := svc.RegisterWorker(context.Background(), ports.RegisterRequest{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestSubmitResult_Success(t *testing.T) {
	notifier := &dummyNotifier{}
	svc := core.NewWorkerGatewayService(notifier)

	err := svc.SubmitResult(context.Background(), ports.ResultRequest{
		JobID:  "job123",
		Status: "COMPLETED",
		Result: "some output",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !notifier.UpdateJobCalled {
		t.Error("expected UpdateJob to be called")
	}
}

func TestSubmitResult_Error(t *testing.T) {
	notifier := &dummyNotifier{ReturnErr: true}
	svc := core.NewWorkerGatewayService(notifier)

	err := svc.SubmitResult(context.Background(), ports.ResultRequest{
		JobID: "job123",
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestHeartbeat_Available_WithJobs(t *testing.T) {
	notifier := &dummyNotifier{}
	svc := core.NewWorkerGatewayService(notifier)

	req := ports.HeartbeatRequest{
		WorkerID: "worker1",
		Status:   "AVAILABLE",
	}

	jobs, err := svc.Heartbeat(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !notifier.UpdateWorkerStatusCalled {
		t.Error("expected UpdateWorkerStatus to be called")
	}
	if !notifier.FetchScheduledJobsCalled {
		t.Error("expected FetchScheduledJobs to be called")
	}
	if len(jobs) != 1 || jobs[0].WorkerID != "worker1" {
		t.Errorf("expected 1 matching job for worker1, got %v", jobs)
	}
}

func TestHeartbeat_Computing(t *testing.T) {
	notifier := &dummyNotifier{}
	svc := core.NewWorkerGatewayService(notifier)

	req := ports.HeartbeatRequest{
		WorkerID: "worker1",
		Status:   "COMPUTING",
	}

	jobs, err := svc.Heartbeat(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(jobs) != 0 {
		t.Errorf("expected no jobs for COMPUTING, got %d", len(jobs))
	}
	if !notifier.UpdateWorkerStatusCalled {
		t.Error("expected UpdateWorkerStatus to be called")
	}
	if notifier.FetchScheduledJobsCalled {
		t.Error("expected FetchScheduledJobs NOT to be called")
	}
}

func TestHeartbeat_Available_ErrorFetchingJobs(t *testing.T) {
	notifier := &dummyNotifier{ReturnErr: true}
	svc := core.NewWorkerGatewayService(notifier)

	req := ports.HeartbeatRequest{
		WorkerID: "worker1",
		Status:   "AVAILABLE",
	}

	jobs, err := svc.Heartbeat(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no fatal error (graceful handling), got %v", err)
	}
	if jobs != nil {
		t.Errorf("expected nil jobs on fetch error, got %v", jobs)
	}
}
