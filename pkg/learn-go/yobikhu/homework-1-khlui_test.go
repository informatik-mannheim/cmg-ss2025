package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestOddEvenHandler(t *testing.T) {
	body := `[29, 8, 3, 4]`
	req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(body))
	w := httptest.NewRecorder()

	// Call handler
	OddEvenHandler(w, req)

	// Check status code
	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", res.Status)
	}

	// Parse actual JSON response
	var actual SortedNumbers
	err := json.Unmarshal(w.Body.Bytes(), &actual)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Define expected result
	expected := SortedNumbers{
		Even: []int{4, 8},
		Odd:  []int{3, 29},
	}

	// Compare values
	if !equalSlices(actual.Even, expected.Even) || !equalSlices(actual.Odd, expected.Odd) {
		t.Errorf("Expected %+v, got %+v", expected, actual)
	}
}

// Helper to compare two int slices
func equalSlices(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestSplitOddEven(t *testing.T) {
	body := []int{1, 2, 3, 4, 5}
	expectedOdd := []int{1, 3, 5}
	expectedEven := []int{2, 4}

	// Call the function to split odd and even numbers
	actualEven, actualOdd := SplitEvenOdd(body)

	// Compare results
	if !equalSlices(actualOdd, expectedOdd) {
		t.Errorf("Expected odd numbers %+v, got %+v", expectedOdd, actualOdd)
	}
	if !equalSlices(actualEven, expectedEven) {
		t.Errorf("Expected even numbers %+v, got %+v", expectedEven, actualEven)
	}
}
