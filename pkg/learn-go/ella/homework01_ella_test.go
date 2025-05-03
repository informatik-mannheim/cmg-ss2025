package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

/* Table-Driven Test covers several edge cases
	including empty Body, !PUT HTTP Method, negative Integers,
	String characters, and invalid content-type*/
func TestEvenOddHandler(t *testing.T) {


	tests := []struct {
		url          string
		name         string
		method       string
		body         string
		code         int
		responseBody string
	}{

		{	url: "/",
			name:         "Valid Request positiv",
			method:       "PUT",
			body:         "[1,2,3,4,5,6]",
			code:         200, // Status OK
			responseBody: `{"Even": [2,4,6], "Odd": [1,3,5]}`,
		},
		{	url: "/",
			name:         "Valid Request negative Integers",
			method:       "PUT",
			body:         "[-1,-2,-3,-4,-5,-6]",
			code:         200, 
			responseBody: `{"Even":[-6,-4,-2],"Odd":[-5,-3,-1]}`,
		},
		{	url: "/",
			name:         "Invalid JSON",
			method:       "PUT",
			body:         "abc123",
			code:         400, // Method not allowed
			responseBody: "Not an Array of Integers",
		},
		{	url: "/",
			name:         "Empty Request",
			method:       "PUT",
			body:         "",
			code:         400,
			responseBody: "Not and Array of Integers",
		},
		{	url: "/",
			name:         "Invalid Method",
			method:       "POST",
			body:         "[1,2,3,4,5,6]",
			code:         400,
			responseBody: "Not a PUT Request",
		},
	}
	// Run subtests for each case
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.url, strings.NewReader(tt.body))

			rr := httptest.NewRecorder()

			handler := http.HandlerFunc(evenOddHandler)
			handler.ServeHTTP(rr, req)

			if tt.code != rr.Code {
				t.Errorf("Expected %d, got %d", tt.code, rr.Code)
			}	
		})
	}
}
/*Tests if the sorting Algorithm 
works accordingly*/
func TestSortNums(t *testing.T) {
	evenOdd := []struct {
		url         string
		name        string
		body        string
		expectedResult []byte
		method      string
		code        int
	}{
		{
			url:         "/",
			name:        "Test even-odd and sort",
			body:        "[1,2,3,4,5,6]",
			expectedResult: []byte(`{"Even":[2,4,6],"Odd":[1,3,5]}`), 
			method:      "PUT",
			code:        200,
		},
	}
	for _, tt := range evenOdd {
		t.Run(tt.name, func(t *testing.T) {

			req := httptest.NewRequest(tt.method, tt.url, strings.NewReader(tt.body))

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