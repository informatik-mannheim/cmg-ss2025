package core_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"io"

	"github.com/informatik-mannheim/cmg-ss2025/services/consumer-gateway/core"
	"github.com/informatik-mannheim/cmg-ss2025/services/consumer-gateway/ports"
	"github.com/informatik-mannheim/cmg-ss2025/services/consumer-gateway/adapters/handler-http"
)




func TestConsumerService(t *testing.T) {
	
	 tests := [] struct {
		name   			string
		url 			string
		method			string
		code   			int
		inputBody		string
		responseBody 	string


	}{
		{
			name: "Forward user login data",
			url: "/auth/login",
			method: "POST",

			inputBody: `{"username" : "Alice Bob", "password": "SuperSecure123"}`,
			responseBody: `{"token": "super-secret-123"}`,
			code: http.StatusOK,

			
		},

		{
			name: "Forward user registration data",
			url: "/auth/register",
			method: "POST",
			code: http.StatusCreated,

			inputBody: `{"username" : "Alice Bob", "password": "SuperSecure123"}`,
			responseBody: `{"token": "super-secret-123"}`,
		},

		{
			name: "Create Job",
			url: "/jobs",
			method: "POST",
			code: http.StatusCreated,
			inputBody: `{"image_id": "123", "location" : "FR", "params":  "abc"}`,

		},
		{
			name: "Get a job result",
			url: "/jobs/{id}/results",
			method: "GET",
			code: http.StatusOK,
		},



	}

		
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := core.NewConsumerService()
			router := handler_http.NewHandler(service)

			body := strings.NewReader(tt.inputBody)
			req := httptest.NewRequest(tt.method, tt.url, body)
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			if rr.Code != tt.code {
				t.Errorf("expected code %d, got %d", tt.code, rr.Code)
				}
	})
}
}