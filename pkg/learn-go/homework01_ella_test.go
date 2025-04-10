package main
 

import (

"net/http"

"net/http/httptest"

"testing"

"strings"

)

  
func TestEvenOddHandlerWithInvalidJSON(t *testing.T) {
	body := "Something"
	req := httptest.NewRequest("PUT", "/", strings.NewReader(body))//Simulates a client-side request
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()//simulates a response and saves it in var w
	evenOddHandler(w, req) //Run function with Server-Response and Client-Request

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 Bad Request, got %d", w.Code)
	}

	expected := "Not an Array of Integers"
	if !strings.Contains(w.Body.String(), expected) {
		t.Errorf("Expected response to contain %q, got %q", expected, w.Body.String())
	}
}
