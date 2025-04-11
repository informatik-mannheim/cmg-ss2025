package main

import (
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"
)

// Unit-Test für separateAndSort
func TestSeparateAndSort(t *testing.T) {
	tests := []struct {
		description         string
		numbers             []int
		expectedEvenNumbers []int
		expectedOddNumbers  []int
	}{
		{
			description:         "Array of mixed numbers",
			numbers:             []int{5, 2, 7, 8, 1, 4},
			expectedEvenNumbers: []int{2, 4, 8},
			expectedOddNumbers:  []int{1, 5, 7},
		},
		{
			description:         "Only even numbers",
			numbers:             []int{4, 2, 8},
			expectedEvenNumbers: []int{2, 4, 8},
			expectedOddNumbers:  []int{},
		},
		{
			description:         "Only odd numbers",
			numbers:             []int{1, 7, 3},
			expectedEvenNumbers: []int{},
			expectedOddNumbers:  []int{1, 3, 7},
		},
		{
			description:         "Empty input",
			numbers:             []int{},
			expectedEvenNumbers: []int{},
			expectedOddNumbers:  []int{},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			actualEvenNumbers, actualOddNumbers := separateAndSort(test.numbers)

			if !slices.Equal(actualEvenNumbers, test.expectedEvenNumbers) {
				t.Errorf("Even slice is wrong. Got %v, but wanted %v", actualEvenNumbers, test.expectedEvenNumbers)
			}

			if !slices.Equal(actualOddNumbers, test.expectedOddNumbers) {
				t.Errorf("Odd slice is wrong. Got %v, but wanted %v", actualOddNumbers, test.expectedOddNumbers)
			}
		})
	}
}

func TestHandler(t *testing.T) {
	tests := []struct {
		description  string
		method       string
		body         string
		status       int
		expectedBody string
	}{
		{
			description:  "Valid PUT request",
			method:       http.MethodPut,
			body:         "[5,2,7,8,1,4]",
			status:       http.StatusOK,
			expectedBody: "Even: [2 4 8], Odd: [1 5 7]",
		},
		{
			description:  "Invalid method",
			method:       http.MethodGet,
			body:         "[5,2,7,8,1,4]",
			status:       http.StatusMethodNotAllowed,
			expectedBody: "Invalid request method\n",
		},
		{
			description:  "Invalid body",
			method:       http.MethodPut,
			body:         "invalid body",
			status:       http.StatusBadRequest,
			expectedBody: "Invalid request body\n",
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			request := httptest.NewRequest(test.method, "/", strings.NewReader(test.body))
			recorder := httptest.NewRecorder()
			handler(recorder, request)
			response := recorder.Result()

			// defer -> Closes response body after the function has finished executing
			defer response.Body.Close()

			if response.StatusCode != test.status {
				t.Errorf("Expected status %d, but got status %d", test.status, response.StatusCode)
			}

			actualBody := recorder.Body.String()

			if actualBody != test.expectedBody {
				t.Errorf("Expected response body %q, got response body %q", test.expectedBody, actualBody)
			}
		})
	}
}
