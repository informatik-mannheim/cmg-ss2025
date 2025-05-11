package auth_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/informatik-mannheim/cmg-ss2025/services/user-management/adapters/auth"
	"github.com/informatik-mannheim/cmg-ss2025/services/user-management/ports"
)

type mockNotifier struct {
	events []string
}

func (m *mockNotifier) Event(msg string) {
	m.events = append(m.events, msg)
}

func (m *mockNotifier) UserRegistered(id string, role string) {}
func (m *mockNotifier) UserLoggedIn(id string)                {}

var _ ports.Notifier = (*mockNotifier)(nil)

func TestAuth0Adapter_RequestTokenFromCredentials(t *testing.T) {
	type fields struct {
		useLive  bool
		notifier *mockNotifier
	}
	type args struct {
		credentials string
	}

	// Setup test HTTP server for live test case
	tokenResponse := map[string]string{"access_token": "live-token-12345", "token_type": "Bearer"}
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tokenResponse)
	}))
	defer mockServer.Close()

	// Set common env vars for live tests
	t.Setenv("AUTH0_TOKEN_URL", mockServer.URL)
	t.Setenv("JWT_AUDIENCE", "https://fake-audience")

	tests := []struct {
		name     string
		fields   fields
		args     args
		setup    func()
		wantPart string
		wantErr  bool
	}{
		{
			name:     "Valid mock credentials",
			fields:   fields{useLive: false, notifier: &mockNotifier{}},
			args:     args{"myclient.mysecret"},
			wantPart: "eyJhbGciOiJS", // mock JWT prefix
			wantErr:  false,
		},
		{
			name:    "Invalid credentials format",
			fields:  fields{useLive: false, notifier: &mockNotifier{}},
			args:    args{"noDotSeparator"},
			wantErr: true,
		},
		{
			name:     "Live token request with mock HTTP server",
			fields:   fields{useLive: true, notifier: &mockNotifier{}},
			args:     args{"testclient.testsecret"},
			wantPart: "live-token-",
			wantErr:  false,
		},
		{
			name:   "Live token request fails on bad URL",
			fields: fields{useLive: true, notifier: &mockNotifier{}},
			args:   args{"client.secret"},
			setup: func() {
				t.Setenv("AUTH0_TOKEN_URL", "http://invalid.invalid")
			},
			wantErr: true,
		},
		{
			name:   "Live token request returns 403",
			fields: fields{useLive: true, notifier: &mockNotifier{}},
			args:   args{"client.secret"},
			setup: func() {
				failServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					http.Error(w, "forbidden", http.StatusForbidden)
				}))
				t.Cleanup(failServer.Close)
				t.Setenv("AUTH0_TOKEN_URL", failServer.URL)
			},
			wantErr: true,
		},
		{
			name:   "Live token request returns invalid JSON",
			fields: fields{useLive: true, notifier: &mockNotifier{}},
			args:   args{"client.secret"},
			setup: func() {
				badJSONServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte("{invalid json"))
				}))
				t.Cleanup(badJSONServer.Close)
				t.Setenv("AUTH0_TOKEN_URL", badJSONServer.URL)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}
			a := auth.New(tt.fields.useLive, tt.fields.notifier)
			got, err := a.RequestTokenFromCredentials(tt.args.credentials)
			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !strings.HasPrefix(got, tt.wantPart) {
				t.Errorf("token = %v, want prefix %v", got, tt.wantPart)
			}
			if len(tt.fields.notifier.events) == 0 {
				t.Errorf("expected at least one notifier event")
			}
		})
	}
}
