// core/service_test.go
package core

import (
	"testing"

	"github.com/informatik-mannheim/cmg-ss2025/services/user-management/model"
)

func TestAddUser(t *testing.T) {
	service := &Service{
		users:    make(map[string]model.User),
		filePath: "",
	}

	id := "test-user"
	role := model.Consumer

	secret, err := service.AddUser(id, role)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if secret == "" {
		t.Error("Expected non-empty secret")
	}

	_, err = service.AddUser(id, role)
	if err == nil {
		t.Error("Expected error for duplicate user, got nil")
	}
}

func TestJobSchedulerSingleton(t *testing.T) {
	service := &Service{
		users:    make(map[string]model.User),
		filePath: "",
	}

	_, err := service.AddUser("scheduler1", model.JobScheduler)
	if err != nil {
		t.Fatalf("Expected first job scheduler to succeed, got: %v", err)
	}

	_, err = service.AddUser("scheduler2", model.JobScheduler)
	if err == nil {
		t.Error("Expected error for second job scheduler, got nil")
	}
}

func TestAuthenticate(t *testing.T) {
	service := &Service{
		users:    make(map[string]model.User),
		filePath: "",
	}

	id := "auth-user"
	role := model.Consumer

	secret, err := service.AddUser(id, role)
	if err != nil {
		t.Fatalf("Failed to add user: %v", err)
	}

	user, err := service.Authenticate(secret)
	if err != nil {
		t.Errorf("Expected authentication to succeed, got error: %v", err)
	}
	if user == nil || user.ID != id {
		t.Errorf("Expected authenticated user with ID %s, got %+v", id, user)
	}

	_, err = service.Authenticate("wrong-secret")
	if err == nil {
		t.Error("Expected authentication to fail with wrong secret")
	}
}
