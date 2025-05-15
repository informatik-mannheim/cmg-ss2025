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

// GetRequest sends an HTTP GET request to the specified URL and decodes the JSON response into the provided type T.
// Returns the decoded data or an error.
func GetRequest[T any](url string) (T, error) {
	// send request
	response, err := http.Get(url)
	if err != nil {
		return *new(T), err
	}

	defer response.Body.Close() // make sure stream is closed

	// check response
	if response.StatusCode != http.StatusOK {
		return *new(T), &RequestError{
			Code:    response.StatusCode,
			Message: fmt.Sprintf("HTTP/GET %s", url),
			Payload: "",
		}
	}

	// decode response
	var data T
	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		return *new(T), err
	}

	return data, nil
}

// PatchRequest sends an HTTP PATCH request with a JSON payload of type T to the specified URL.
// The response is decoded into type R. Returns the decoded response or an error.
func PatchRequest[T any, R any](url string, payload T) (R, error) {
	// prepare request
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return *new(R), err
	}

	request, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return *new(R), err
	}
	request.Header.Set("Content-Type", "application/json")

	// send request
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return *new(R), err
	}
	defer response.Body.Close() // make sure stream is closed

	// check response
	if response.StatusCode != http.StatusOK {
		return *new(R), &RequestError{
			Code:    response.StatusCode,
			Message: fmt.Sprintf("HTTP/PATCH %s", url),
			Payload: string(jsonPayload),
		}
	}

	// decode response
	var data R
	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		return *new(R), err
	}

	return data, nil
}

// PutRequest sends an HTTP PUT request with a JSON payload of type T to the specified URL.
// The response is decoded into type R. Returns the decoded response or an error.
func PutRequest[T any, R any](url string, payload T) (R, error) {
	// prepare request
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return *new(R), err
	}

	request, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return *new(R), err
	}
	request.Header.Set("Content-Type", "application/json")

	// send request
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return *new(R), err
	}
	defer response.Body.Close() // make sure stream is closed

	// check response
	if response.StatusCode != http.StatusOK {
		return *new(R), &RequestError{
			Code:    response.StatusCode,
			Message: fmt.Sprintf("HTTP/PUT %s", url),
			Payload: string(jsonPayload),
		}
	}

	// decode response
	var data R
	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		return *new(R), err
	}

	return data, nil
}
