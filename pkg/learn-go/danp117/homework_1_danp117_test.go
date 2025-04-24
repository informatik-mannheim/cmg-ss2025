package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		body         string
		wantStatus   int
		wantResponse string
	}{
		{
			name:         "Valid input",
			method:       http.MethodPut,
			body:         `[29, 8, 3, 4]`,
			wantStatus:   http.StatusOK,
			wantResponse: "even : [4, 8], odd : [3, 29]",
		},
		{
			name:         "Invalid JSON",
			method:       http.MethodPut,
			body:         `[1, 2, "bad"]`,
			wantStatus:   http.StatusBadRequest,
			wantResponse: "Invalid JSON input",
		},
		{
			name:         "Wrong method",
			method:       http.MethodGet,
			body:         `[1, 2, 3]`,
			wantStatus:   http.StatusMethodNotAllowed,
			wantResponse: "This method is not allowed",
		},
		{
			name:         "Empty array",
			method:       http.MethodPut,
			body:         `[]`,
			wantStatus:   http.StatusOK,
			wantResponse: "even : [], odd : []",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler(w, req)

			res := w.Result()
			if res.StatusCode != tt.wantStatus {
				t.Errorf("Expected status %d, got %d", tt.wantStatus, res.StatusCode)
			}

			gotBody := strings.TrimSpace(w.Body.String())
			if !strings.HasPrefix(gotBody, tt.wantResponse) {
				t.Errorf("Expected body to start with %q, got %q", tt.wantResponse, gotBody)
			}
		})
	}
}
