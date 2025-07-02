package core_test

import (
	"context"
	"errors"
	"testing"

	"github.com/informatik-mannheim/cmg-ss2025/pkg/logging"
	"github.com/informatik-mannheim/cmg-ss2025/services/worker-gateway/core"
	"github.com/informatik-mannheim/cmg-ss2025/services/worker-gateway/ports"
)

func init() {
	logging.Init("Worker Registry")
}

// --- Dummy RegistryService für Tests ---
type dummyRegistryService struct {
	RegisterWorkerCalled     bool
	UpdateWorkerStatusCalled bool
	ReturnErr                bool
}

func (d *dummyRegistryService) RegisterWorker(ctx context.Context, req ports.RegisterRequest) (*ports.RegisterRespose, error) {
	d.RegisterWorkerCalled = true
	if d.ReturnErr {
		return nil, errors.New("register worker error")
	}
	return &ports.RegisterRespose{
		ID:     "worker123",
		Status: "AVAILABLE",
		Zone:   "DE",
	}, nil
}

func (d *dummyRegistryService) UpdateWorkerStatus(ctx context.Context, req ports.HeartbeatRequest) error {
	d.UpdateWorkerStatusCalled = true
	if d.ReturnErr {
		return errors.New("update worker status error")
	}
	return nil
}

// --- Dummy JobService für Tests ---
type dummyJobService struct {
	UpdateJobCalled          bool
	FetchScheduledJobsCalled bool
	ReturnErr                bool
}

func (d *dummyJobService) UpdateJob(ctx context.Context, req ports.ResultRequest) error {
	d.UpdateJobCalled = true
	if d.ReturnErr {
		return errors.New("update job error")
	}
	return nil
}

func (d *dummyJobService) FetchScheduledJobs(ctx context.Context) ([]ports.Job, error) {
	d.FetchScheduledJobsCalled = true
	if d.ReturnErr {
		return nil, errors.New("fetch jobs error")
	}
	return []ports.Job{
		{ID: "job1", WorkerID: "worker1", Status: "SCHEDULED"},
		{ID: "job2", WorkerID: "worker2", Status: "SCHEDULED"},
	}, nil
}

// --- Dummy UserClient für Tests ---
type dummyUserClient struct {
	GetTokenCalled bool
	ReturnErr      bool
}

func (d *dummyUserClient) GetToken(ctx context.Context) (string, error) {
	d.GetTokenCalled = true
	if d.ReturnErr {
		return "", errors.New("user token error")
	}
	return "mocked.secret.token", nil
}

// --- Helper zum Erstellen des Services ---
func newTestWorkerGatewayService(reg ports.RegistryService, job ports.JobService, user ports.UserService) *core.WorkerGatewayService {
	return core.NewWorkerGatewayService(reg, job, user)
}

// --- Tests ---

func TestRegisterWorker_WithUserClient_Success(t *testing.T) {
	reg := &dummyRegistryService{}
	job := &dummyJobService{}
	user := &dummyUserClient{}
	svc := newTestWorkerGatewayService(reg, job, user)

	req := ports.RegisterRequest{
		Zone: "EU",
	}

	resp, err := svc.Register(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !user.GetTokenCalled {
		t.Error("expected RegisterProvider to be called")
	}
	if !reg.RegisterWorkerCalled {
		t.Error("expected RegisterWorker to be called")
	}
	if resp == nil || resp.ID != "worker123" {
		t.Errorf("unexpected register response: %+v", resp)
	}
}

func TestRegisterWorker_UserClientFails(t *testing.T) {
	reg := &dummyRegistryService{}
	job := &dummyJobService{}
	user := &dummyUserClient{ReturnErr: true}
	svc := newTestWorkerGatewayService(reg, job, user)

	_, err := svc.Register(context.Background(), ports.RegisterRequest{Zone: "EU"})
	if err == nil {
		t.Fatal("expected error due to user client failure, got nil")
	}
}

func TestRegisterWorker_RegistryFails(t *testing.T) {
	reg := &dummyRegistryService{ReturnErr: true}
	job := &dummyJobService{}
	user := &dummyUserClient{}
	svc := newTestWorkerGatewayService(reg, job, user)

	_, err := svc.Register(context.Background(), ports.RegisterRequest{Zone: "EU"})
	if err == nil {
		t.Fatal("expected error due to registry failure, got nil")
	}
}

func TestSubmitResult_Success(t *testing.T) {
	reg := &dummyRegistryService{}
	job := &dummyJobService{}
	user := &dummyUserClient{}
	svc := newTestWorkerGatewayService(reg, job, user)

	err := svc.Result(context.Background(), ports.ResultRequest{
		JobID:  "job123",
		Status: "COMPLETED",
		Result: "some output",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !job.UpdateJobCalled {
		t.Error("expected UpdateJob to be called")
	}
}

func TestSubmitResult_Error(t *testing.T) {
	reg := &dummyRegistryService{}
	job := &dummyJobService{ReturnErr: true}
	user := &dummyUserClient{}
	svc := newTestWorkerGatewayService(reg, job, user)

	err := svc.Result(context.Background(), ports.ResultRequest{
		JobID: "job123",
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestHeartbeat_Available_WithJobs(t *testing.T) {
	reg := &dummyRegistryService{}
	job := &dummyJobService{}
	user := &dummyUserClient{}
	svc := newTestWorkerGatewayService(reg, job, user)

	req := ports.HeartbeatRequest{
		WorkerID: "worker1",
		Status:   "AVAILABLE",
	}

	jobs, err := svc.Heartbeat(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !reg.UpdateWorkerStatusCalled {
		t.Error("expected UpdateWorkerStatus to be called")
	}
	if !job.FetchScheduledJobsCalled {
		t.Error("expected FetchScheduledJobs to be called")
	}
	if len(jobs) != 1 || jobs[0].WorkerID != "worker1" {
		t.Errorf("expected 1 matching job for worker1, got %v", jobs)
	}
}

func TestHeartbeat_Computing(t *testing.T) {
	reg := &dummyRegistryService{}
	job := &dummyJobService{}
	user := &dummyUserClient{}
	svc := newTestWorkerGatewayService(reg, job, user)

	req := ports.HeartbeatRequest{
		WorkerID: "worker1",
		Status:   "RUNNING",
	}

	jobs, err := svc.Heartbeat(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(jobs) != 0 {
		t.Errorf("expected no jobs for RUNNING, got %d", len(jobs))
	}
	if !reg.UpdateWorkerStatusCalled {
		t.Error("expected UpdateWorkerStatus to be called")
	}
	if job.FetchScheduledJobsCalled {
		t.Error("expected FetchScheduledJobs NOT to be called")
	}
}

func TestHeartbeat_Available_ErrorFetchingJobs(t *testing.T) {
	reg := &dummyRegistryService{}
	job := &dummyJobService{ReturnErr: true}
	user := &dummyUserClient{}
	svc := newTestWorkerGatewayService(reg, job, user)

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
