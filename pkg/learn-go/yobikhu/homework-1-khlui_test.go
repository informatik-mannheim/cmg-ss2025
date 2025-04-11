package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestOddEvenHandler(t *testing.T) {
	// Test case: Valid input
	body := `[29, 8, 3, 4]`

	// Create a new request with the PUT method and the test body
	req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(body))
	w := httptest.NewRecorder()

	// Call the handler function
	OddEvenHandler(w, req)

	// Check the response status code
	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", res.Status)
	}

	// Check the response body
	// Remove leading and trailing whitespace from the actual response body
	// and compare it with the expected body
	expectedBody := "even: [4,8], odd: [3,29]"
	actual := strings.TrimSpace(w.Body.String())
	if actual != expectedBody {
		t.Errorf("Expected body %q, got %q", expectedBody, actual)
	}

}
