package main

import (
	"net/http"
	"net/httpt/httptest"
	"testing"
	"strings"
)

func TestEvenOddHandlerInvalidJSON(t *testing.T) {
	body := "something" //Input thats not an Integer Array
	req := httptest.NewRequest("PUT", "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	evenOddHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 Bad Request, got %d", w.Code)
	}

	expected := "Not an Array of Integers"
	if !strings.Contains(w.Body.String(), expected) {
		t.Errorf("Expected response to contain %q, got %q", expected, w.Body.String())
	}
}

func TestEvenOddHandlerValidJSON(t *testing.T) {
	body := "[1,2,3,4,5,6]" //Input thats not an Integer Array
	req := httptest.NewRequest("PUT", "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	evenOddHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 Bad Request, got %d", w.Code)
	}

	expected := "An Array of Integers"
	if !strings.Contains(w.Body.String(), expected) {
		t.Errorf("Expected response to contain Integers", expected, w.Body.String())
	}
}
