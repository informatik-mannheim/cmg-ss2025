package core

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"worker-daemon/internal/config"
	"worker-daemon/internal/ports"
)

// DummyWorkerGateway simuliert das API-Interface
type DummyWorkerGateway struct {
	RegisterCalled      bool
	SendHeartbeatCalled bool
	SendResultCalled    bool

	// Simuliere Fehler
	RegisterErr      error
	SendHeartbeatErr error
	SendResultErr    error

	JobsToReturn []ports.Job
	ReceivedJobs []ports.Job
}

func (d *DummyWorkerGateway) Register(key, zone string) (*ports.RegisterResponse, error) {
	d.RegisterCalled = true
	if d.RegisterErr != nil {
		return &ports.RegisterResponse{}, d.RegisterErr
	}
	return &ports.RegisterResponse{
		ID:    "worker123",
		Token: "token123",
	}, nil
}

func (d *DummyWorkerGateway) SendHeartbeat(workerID, status, token string) ([]ports.Job, error) {
	d.SendHeartbeatCalled = true
	if d.SendHeartbeatErr != nil {
		return nil, d.SendHeartbeatErr
	}
	return d.JobsToReturn, nil
}

func (d *DummyWorkerGateway) SendResult(job ports.Job, token string) error {
	d.SendResultCalled = true
	d.ReceivedJobs = append(d.ReceivedJobs, job)
	return d.SendResultErr
}

func TestDaemon_HeartbeatLoop_RegisterFails(t *testing.T) {
	dummyAPI := DummyWorkerGateway{
		RegisterErr: errors.New("register failed"),
	}
	cfg := config.Config{
		Secret:                   "key",
		Zone:                     "zone",
		HeartbeatIntervalSeconds: 1,
	}
	d := NewDaemon(cfg, &dummyAPI)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Starte HeartbeatLoop, der Registrierung schlägt fehl und der Loop soll abbrechen
	d.StartHeartbeatLoop(ctx)

	if !dummyAPI.RegisterCalled {
		t.Error("expected Register to be called")
	}
}

func TestDaemon_HeartbeatLoop_ProcessJob(t *testing.T) {
	dummyAPI := &DummyWorkerGateway{
		JobsToReturn: []ports.Job{
			{
				ID: "job1",
				Image: ports.ContainerImage{
					Name:    "alpine",
					Version: "latest",
				},
				AdjustmentParameters: map[string]string{
					"echo":  "hello",
					"param": "value",
				},
			},
		},
	}
	cfg := config.Config{
		Secret:                   "key",
		Zone:                     "zone",
		HeartbeatIntervalSeconds: 1,
	}
	d := NewDaemon(cfg, dummyAPI)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Starte Heartbeat Loop in goroutine, um nicht zu blockieren
	go d.StartHeartbeatLoop(ctx)

	// Warte eine Weile, damit Loop mindestens einmal ausgeführt wird
	time.Sleep(3 * time.Second)

	if !dummyAPI.RegisterCalled {
		t.Error("expected Register to be called")
	}

	if !dummyAPI.SendHeartbeatCalled {
		t.Error("expected SendHeartbeat to be called")
	}

	// Da Jobs zurückgegeben wurden, sollte auch SendResult mindestens einmal aufgerufen werden
	if !dummyAPI.SendResultCalled {
		t.Error("expected SendResult to be called")
	}

	if len(dummyAPI.ReceivedJobs) == 0 {
		t.Error("expected at least one job to be sent as result")
	}

	// Prüfe, ob der Jobstatus auf "DONE" oder "ERROR" gesetzt wurde
	job := dummyAPI.ReceivedJobs[0]
	if job.Status != "DONE" && job.Status != "ERROR" {
		t.Errorf("expected job status DONE or ERROR, got %s", job.Status)
	}
}

func TestDaemon_HeartbeatLoop_HeartbeatFails(t *testing.T) {
	dummyAPI := &DummyWorkerGateway{
		JobsToReturn:     nil,
		SendHeartbeatErr: errors.New("heartbeat failed"),
		RegisterErr:      nil,
		SendResultErr:    nil,
	}
	cfg := config.Config{
		Secret:                   "key",
		Zone:                     "zone",
		HeartbeatIntervalSeconds: 1,
	}
	d := NewDaemon(cfg, dummyAPI)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	go d.StartHeartbeatLoop(ctx)

	time.Sleep(2 * time.Second)

	if !dummyAPI.RegisterCalled {
		t.Error("expected Register to be called")
	}

	if !dummyAPI.SendHeartbeatCalled {
		t.Error("expected SendHeartbeat to be called")
	}

	// Da Heartbeat Fehler liefert, sollte SendResult NICHT aufgerufen werden
	if dummyAPI.SendResultCalled {
		t.Error("did not expect SendResult to be called when Heartbeat fails")
	}
}

func TestComputeJob_Success(t *testing.T) {
	job := ports.Job{
		ID: "test-123",
		Image: ports.ContainerImage{
			Name:    "alpine",
			Version: "latest",
		},
		AdjustmentParameters: map[string]string{
			"echo": "hello test",
		},
	}

	fmt.Printf("Starte Job mit Image: %s:%s\n", job.Image.Name, job.Image.Version)
	fmt.Printf("AdjustmentParameters: %v\n", job.AdjustmentParameters)

	result := computeJob(job)

	fmt.Println("Job abgeschlossen. Ergebnis:")
	fmt.Printf("Status:       %s\n", result.Status)
	fmt.Printf("Result:       %q\n", result.Result)
	fmt.Printf("ErrorMessage: %q\n", result.ErrorMessage)

	if result.Status != "DONE" {
		t.Errorf("Expected status DONE, got %s", result.Status)
	}

	if !strings.Contains(result.Result, "hello test") {
		t.Errorf("Expected result to contain 'hello test', got: %s", result.Result)
	}

	if result.ErrorMessage != "" {
		t.Errorf("Expected no error message, got: %s", result.ErrorMessage)
	}
}

func TestComputeJob_InvalidImage(t *testing.T) {
	job := ports.Job{
		ID: "fail-999",
		Image: ports.ContainerImage{
			Name:    "nonexistent-image-akjbfsfkjbdgnsdfkjv",
			Version: "never",
		},
		AdjustmentParameters: map[string]string{
			"echo": "fail",
		},
	}

	fmt.Printf("Starte Job mit Image: %s:%s\n", job.Image.Name, job.Image.Version)
	fmt.Printf("AdjustmentParameters: %v\n", job.AdjustmentParameters)

	result := computeJob(job)

	fmt.Println("Job abgeschlossen. Ergebnis:")
	fmt.Printf("Status:       %s\n", result.Status)
	fmt.Printf("Result:       %q\n", result.Result)
	fmt.Printf("ErrorMessage: %q\n", result.ErrorMessage)

	if result.Status != "ERROR" {
		t.Errorf("Expected status ERROR, got %s", result.Status)
	}

	if result.Result != "" {
		t.Errorf("Expected empty result, got: %s", result.Result)
	}

	if result.ErrorMessage == "" {
		t.Errorf("Expected error message, got empty")
	}
}
