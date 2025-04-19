package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
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
			code:         400, // Method not allowed
			responseBody: "Not an Array of Integers",
		},
		{url: "/",
			name:         "Empty Request",
			contentType:  "application/json",
			method:       "PUT",
			body:         " ",
			code:         400,
			responseBody: "Not and Array of Integers",
		},
		{url: "/",
			name:         "Invalid Method",
			contentType:  "application/json",
			method:       "POST",
			body:         "[1,2,3,4,5,6]",
			code:         400,
			responseBody: "Not a PUT Request",
		},
		{url: "/",
			name:         "Invalid Content-Type",
			contentType:  "image/gif",
			method:       "PUT",
			body:         "[1,2,3,45,6]",
			code:         400,
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

			if tt.code != rr.Code {
				t.Errorf("Expected %d, got %d", tt.code, rr.Code)
			}

		})
	}
}

func TestSortNums(t *testing.T) {
	evenOdd := []struct {
		url         string
		name        string
		body        string
		expectedResult []byte
		method      string
		contentType string
		code        int
	}{
		{
			url:         "/",
			name:        "Test even-odd and sort",
			body:        "[1,2,3,4,5,6]",
			expectedResult: []byte(`{"Even":[2,4,6],"Odd":[1,3,5]}`), 
			method:      "PUT",
			contentType: "application/json",
			code:        200,
		},
	}
	for _, tt := range evenOdd {
		t.Run(tt.name, func(t *testing.T) {

			req := httptest.NewRequest(tt.method, tt.url, strings.NewReader(tt.body))
			req.Header.Set("Content-Type", tt.contentType)

			rr := httptest.NewRecorder()

			handler := http.HandlerFunc(evenOddHandler)
			handler.ServeHTTP(rr, req)

			if tt.code != rr.Code {
				t.Errorf("Expected %d, got %d", tt.code, rr.Code)

			} else if !bytes.Equal(tt.expectedResult, rr.Body.Bytes()){
				t.Errorf("Expected body %s, got %s", tt.expectedResult, rr.Body.Bytes())
			}
				
		})
	}
}
