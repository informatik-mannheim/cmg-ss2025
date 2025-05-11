package handler_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"crypto/sha256"
	"encoding/hex"

	"github.com/informatik-mannheim/cmg-ss2025/services/user-management/adapters/handler-http"
	"github.com/informatik-mannheim/cmg-ss2025/services/user-management/ports"
)

type mockNotifier struct {
	events []string
	logins []string
	regs   []string
}

func (m *mockNotifier) Event(msg string, _ context.Context) {
	m.events = append(m.events, msg)
}

func (m *mockNotifier) UserLoggedIn(id string, _ context.Context) {
	m.logins = append(m.logins, id)
}

func (m *mockNotifier) UserRegistered(id, role string, _ context.Context) {
	m.regs = append(m.regs, id+":"+role)
}

type dummyAuth struct{}

func (dummyAuth) RequestTokenFromCredentials(credentials string) (string, error) {
	return "dummytoken", nil
}

func TestHTTPHandler_RegisterHandler(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		isAdmin    bool
		wantStatus int
		eventCheck string
		regsCount  int
	}{
		{
			name:       "Valid consumer registration",
			body:       `{"role":"consumer"}`,
			isAdmin:    true,
			wantStatus: http.StatusCreated,
			eventCheck: "Registration successful for role: consumer",
			regsCount:  1,
		},
		{
			name:       "Unauthorized register attempt",
			body:       `{"role":"consumer"}`,
			isAdmin:    false,
			wantStatus: http.StatusUnauthorized,
			eventCheck: "Unauthorized register attempt",
			regsCount:  0,
		},
		{
			name:       "Invalid JSON payload",
			body:       `{invalid`,
			isAdmin:    true,
			wantStatus: http.StatusBadRequest,
			eventCheck: "Invalid register request payload",
			regsCount:  0,
		},
		{
			name:       "Unauthorized job scheduler",
			body:       `{"role":"job scheduler"}`,
			isAdmin:    false,
			wantStatus: http.StatusForbidden,
			eventCheck: "Unauthorized attempt to create Job Scheduler",
			regsCount:  0,
		},
		{
			name:       "Invalid role",
			body:       `{"role":"hacker"}`,
			isAdmin:    true,
			wantStatus: http.StatusBadRequest,
			eventCheck: "Invalid register request payload",
			regsCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &mockNotifier{}
			h := handler.New(dummyAuth{}, false, func(string) bool { return tt.isAdmin }, func() ports.Notifier { return n })

			req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer([]byte(tt.body)))
			w := httptest.NewRecorder()
			h.RegisterHandler(w, req)

			res := w.Result()
			if res.StatusCode != tt.wantStatus {
				t.Errorf("got status %v, want %v", res.StatusCode, tt.wantStatus)
			}
			found := false
			for _, e := range n.events {
				if strings.Contains(e, tt.eventCheck) {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("expected event containing: %q", tt.eventCheck)
			}
			if len(n.regs) != tt.regsCount {
				t.Errorf("got %d registrations, want %d", len(n.regs), tt.regsCount)
			}
		})
	}
}

func TestHTTPHandler_LoginHandler(t *testing.T) {
	tests := []struct {
		name        string
		body        string
		wantStatus  int
		eventCheck  string
		loginsCount int
	}{
		{
			name:        "Valid login request",
			body:        `{"secret":"id.secret"}`,
			wantStatus:  http.StatusOK,
			eventCheck:  "Login successful for client: id",
			loginsCount: 1,
		},
		{
			name:        "Invalid login JSON",
			body:        `notjson`,
			wantStatus:  http.StatusBadRequest,
			eventCheck:  "Invalid login request format",
			loginsCount: 0,
		},
		{
			name:        "Malformed secret",
			body:        `{"secret":"missingdot"}`,
			wantStatus:  http.StatusOK,
			eventCheck:  "Login successful for client: ",
			loginsCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &mockNotifier{}
			h := handler.New(dummyAuth{}, false, func(string) bool { return true }, func() ports.Notifier { return n })

			req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer([]byte(tt.body)))
			w := httptest.NewRecorder()
			h.LoginHandler(w, req)

			res := w.Result()
			if res.StatusCode != tt.wantStatus {
				t.Errorf("got status %v, want %v", res.StatusCode, tt.wantStatus)
			}
			found := false
			for _, e := range n.events {
				if strings.Contains(e, tt.eventCheck) {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("expected event containing: %q", tt.eventCheck)
			}
			if len(n.logins) != tt.loginsCount {
				t.Errorf("got %d logins, want %d", len(n.logins), tt.loginsCount)
			}
		})
	}
}

func TestIsAdmin(t *testing.T) {
	secret := "adminsecret"
	hash := sha256.Sum256([]byte(secret))
	expected := hex.EncodeToString(hash[:])
	os.Setenv("ADMIN_SECRET_HASH", expected)

	if !handler.IsAdmin(secret) {
		t.Error("IsAdmin() returned false for valid admin secret")
	}
	if handler.IsAdmin("wrong") {
		t.Error("IsAdmin() returned true for invalid admin secret")
	}
}
