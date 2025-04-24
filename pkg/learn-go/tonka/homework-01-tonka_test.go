package tonka

import (
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"
)

func TestSeperateAndSort(t *testing.T) {
	var tests = []struct {
		description  string
		input        []int
		expectedEven []int
		expectedOdd  []int
	}{
		{
			description:  "Test challenging sorting and seperation",
			input:        []int{6, 5, 4, 1, 2, 3},
			expectedEven: []int{2, 4, 6},
			expectedOdd:  []int{1, 3, 5},
		},
		{
			description:  "Test empty input",
			input:        []int{},
			expectedEven: []int{},
			expectedOdd:  []int{},
		},
		{
			description:  "Test all even numbers",
			input:        []int{2, 4, 6},
			expectedEven: []int{2, 4, 6},
			expectedOdd:  []int{},
		},
		{
			description:  "Test all odd numbers",
			input:        []int{1, 3, 5},
			expectedEven: []int{},
			expectedOdd:  []int{1, 3, 5},
		},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			actual := seperateAndSort(test.input)
			if !slices.Equal(actual.Even, test.expectedEven) {
				t.Errorf("Expected even numbers %v, got %v", test.expectedEven, actual.Even)
			}
			if !slices.Equal(actual.Odd, test.expectedOdd) {
				t.Errorf("Expected odd numbers %v, got %v", test.expectedOdd, actual.Odd)
			}
		})
	}
}

func TestHandler(t *testing.T) {
	var tests = []struct {
		description string
		method      string
		input       string
		expected    string
		statusCode  int
	}{
		{
			description: "Test valid input",
			method:      http.MethodPut,
			input:       `[6,5,4,1,2,3]`,
			expected:    `{"even":[2,4,6],"odd":[1,3,5]}`,
			statusCode:  http.StatusOK,
		},
		{
			description: "Test empty input",
			method:      http.MethodPut,
			input:       `[]`,
			expected:    `{"even":[],"odd":[]}`,
			statusCode:  http.StatusOK,
		},
		{
			description: "Test invalid input",
			method:      http.MethodPut,
			input:       `[6,5,4,1,2,"a"]`,
			expected:    `Invalid request body`,
			statusCode:  http.StatusBadRequest,
		},
		{
			description: "Test invalid method",
			method:      http.MethodGet,
			input:       `[6,5,4,1,2,3]`,
			expected:    `Invalid request method`,
			statusCode:  http.StatusMethodNotAllowed,
		},
		{
			description: "Test invalid request body",
			method:      http.MethodPut,
			input:       `invalid json`,
			expected:    `Invalid request body`,
			statusCode:  http.StatusBadRequest,
		},
		{
			description: "Test curl of slides",
			method:      http.MethodPut,
			input:       `[29,8,3,4]`,
			expected:    `{"even":[4,8],"odd":[3,29]}`,
			statusCode:  http.StatusOK,
		},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			request := httptest.NewRequest(test.method, "/", strings.NewReader(test.input))
			response := httptest.NewRecorder()
			handler(response, request)
			if response.Code != test.statusCode {
				t.Errorf("Expected status code %d, got %d", test.statusCode, response.Code)
			}
			actualBody := response.Body.String()
			if strings.TrimRight(actualBody, "\n") != test.expected {
				t.Errorf("Expected body %s, got %s", test.expected, actualBody)
			}
		})
	}
}
