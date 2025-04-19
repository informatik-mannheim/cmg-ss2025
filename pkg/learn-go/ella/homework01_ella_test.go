package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"golang.org/x/oauth2/odnoklassniki"
	"google.golang.org/api/eventarc/v1"
)

func TestEvenOddHandler(t *testing.T) {

	// Table-Driven Test
	tests := []struct {
		url          string
		name         string
		contentType  string
		method       string
		body         string
		code         int
		responseBody string
	}{

		{url: "/",
			name:         "Valid Request",
			contentType:  "application/json",
			method:       "PUT",
			body:         "[1,2,3,4,5,6]",
			code:         200, // Status OK
			responseBody: `{"Even": [2,4,6], "Odd": [1,3,5]}`,
		},
		{url: "/",
			name:         "Invalid JSON",
			contentType:  "application/json",
			method:       "PUT",
			body:         "abc123",
			code:         405, // Method not allowed
			responseBody: "Not an Array of Integers",
		},
		{url: "/",
			name:         "Invalid Method",
			contentType:  "application/json",
			method:       "POST",
			body:         "[1,2,3,4,5,6]",
			code:         405,
			responseBody: "Not a PUT Request",
		},
		{url: "/",
			name:         "Invalid Content-Type",
			contentType:  "image/gif",
			method:       "PUT",
			body:         "[1,2,3,45,6]",
			code:         405,
			responseBody: "Not in JSON Format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.url, strings.NewReader(tt.body))
			req.Header.Set("Content-Type", tt.contentType)

			rr := httptest.NewRecorder()

			handler := http.HandlerFunc(evenOddHandler)
			handler.ServeHTTP(rr, req)

		})
	}
}

func TestSortNums(t *testing.T){
	tests := []struct {
		string name
		int input
		string result
	}
 
		{ name: "Test even-odd and sort"
		input: "[1,2,3,4,5,6]"
		result: "{"Even": "2,4,6", "Odd": "1,3,5"}"
		},

	}

}
