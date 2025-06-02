package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// RequestError represents an error returned from an HTTP request, including status code, message, and payload.
type RequestError struct {
	Code    int
	Message string
	Payload string
}

var _ error = (*RequestError)(nil)

// Error implements the error interface for RequestError.
func (e *RequestError) Error() string {
	return fmt.Sprintf("Request failed with status code %d with payload <%s>: %s", e.Code, e.Payload, e.Message)
}

func isNotStatusCodeSuccess(statusCode int) bool {
	// Check if the status code is in the range of 200-299
	return !(statusCode >= 200 && statusCode < 300)
}

// doRequest is a helper to send an HTTP request, check status, and decode the response.
func doRequest[Resp any](req *http.Request, jsonPayload []byte, method, url string) (Resp, int, error) {
	client := &http.Client{}
	response, err := client.Do(req)

	if err != nil {
		return *new(Resp), -1, err
	}
	defer response.Body.Close()

	if isNotStatusCodeSuccess(response.StatusCode) {
		return *new(Resp), response.StatusCode, &RequestError{
			Code:    response.StatusCode,
			Message: fmt.Sprintf("HTTP/%s %s", method, url),
			Payload: string(jsonPayload),
		}
	}

	var data Resp
	if response.ContentLength == 0 {
		return data, response.StatusCode, nil
	}

	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		return *new(Resp), response.StatusCode, err
	}

	return data, response.StatusCode, nil
}

// GetRequest sends an HTTP GET request to the specified URL and decodes the JSON response into the provided type T.
// Returns the decoded data, a status code or an error.
// If the status code -1 is returned, it indicates an error occurred during the request.
// If the status code is not in the range of 200-299, it returns a RequestError with the status code and message.
func GetRequest[T any](url string) (T, int, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return *new(T), -1, err
	}

	return doRequest[T](req, nil, http.MethodGet, url)
}

// PatchRequest sends an HTTP PATCH request with a JSON payload of type T to the specified URL.
// The response is decoded into type R. Returns the decoded response or an error.
// If the status code -1 is returned, it indicates an error occurred during the request.
// If the status code is not in the range of 200-299, it returns a RequestError with the status code and message.
func PatchRequest[T any, R any](url string, payload T) (R, int, error) {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return *new(R), -1, err
	}

	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return *new(R), -1, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Close = true

	return doRequest[R](req, jsonPayload, http.MethodPatch, url)
}

// PutRequest sends an HTTP PUT request with a JSON payload of type T to the specified URL.
// The response is decoded into type R. Returns the decoded response or an error.
// If the status code -1 is returned, it indicates an error occurred during the request.
// If the status code is not in the range of 200-299, it returns a RequestError with the status code and message.
func PutRequest[T any, R any](url string, payload T) (R, int, error) {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return *new(R), -1, err
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return *new(R), -1, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Close = true

	return doRequest[R](req, jsonPayload, http.MethodPut, url)
}
