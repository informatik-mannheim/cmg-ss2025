package tonka_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"
	"tonka"
)

func TestSeperateAndSort(t *testing.T) {
	var numbers = []int{6, 5, 4, 1, 2, 3}
	var expectedEven = []int{2, 4, 6}
	var expectedOdd = []int{1, 3, 5}

	var actualSortedNumbers = tonka.SeperateAndSort(numbers)

	// Check if the evens and odds are equal the expected values
	if !slices.Equal(actualSortedNumbers.Even, expectedEven) {
		t.Errorf("Expected even numbers: %v, got: %v", expectedEven, actualSortedNumbers.Even)
	}
	if !slices.Equal(actualSortedNumbers.Odd, expectedOdd) {
		t.Errorf("Expected odd numbers: %v, got: %v", expectedOdd, actualSortedNumbers.Odd)
	}
}

func TestHandler(t *testing.T) {
	var body io.Reader = strings.NewReader(`[6,5,4,1,2,3]`)
	request := httptest.NewRequest(http.MethodPut, "/", body)
	response := httptest.NewRecorder()

	tonka.Handler(response, request)

	if response.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.Code)
	}

	expectedBody := `{"even":[2,4,6],"odd":[1,3,5]}`
	actualBody := response.Body.String()
	if strings.TrimRight(actualBody, "\n") != expectedBody {
		t.Errorf("Expected body %s, got %s", expectedBody, actualBody)
	}
}
