package core_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

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
		responseBody 	string

		zone 			string
		params 			string
		body 			string

		username string
		password string
		token	 string

	}{
		{
			name: "Forward user login data",
			url: "/auth/login",
			method: "POST",
			code: http.StatusOK,

			username: "Alice Bob",
			password: "SuperSecurePassword123",
			token: "super-secret-123",
		},

		{
			name: "Forward user registration data",
			url: "/auth/register",
			method: "POST",
			code: http.StatusCreated,

			username: "Alice Bob",
			password: "SuperSecurePassword123",
			token: "super-secret-123",
		},

		{
			name: "Create Job",
			url: "/jobs",
			method: "POST",
			code: http.StatusCreated,
			zone: "GER",
			params: "-abc",

		},

		{
			name: "Get a job result",
			url: "/jobs/{id}/results",
			method: "GET",
			code: http.StatusOK,
		},


	}
	
		{
		
		for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			
			if (tt.name == "Create Job") { 

				service := core.NewConsumerService()
				router := handler_http.NewHandler(service)

				req := httptest.NewRequest("POST", "/jobs", body)
				rr := httptest.NewRecorder()

				router.ServeHTTP(rr, req)

		
		})
	} }

		
}
