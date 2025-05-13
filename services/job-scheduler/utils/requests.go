package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type RequestError struct {
	Code    int
	Message string
	Payload string
}

func (e *RequestError) Error() string {
	return fmt.Sprintf("Request failed with status code %d with payload <%s>: %s", e.Code, e.Payload, e.Message)
}

func GetRequest[T any](url string) (T, error) {
	response, err := http.Get(url)
	if err != nil {
		return *new(T), err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return *new(T), &RequestError{
			Code:    response.StatusCode,
			Message: fmt.Sprintf("HTTP/GET %s", url),
			Payload: "",
		}
	}

	var data T
	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		return *new(T), err
	}

	return data, nil
}

func PatchRequest[T any, R any](url string, payload T) (R, error) {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return *new(R), err
	}

	request, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return *new(R), err
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return *new(R), err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return *new(R), &RequestError{
			Code:    response.StatusCode,
			Message: fmt.Sprintf("HTTP/PATCH %s", url),
			Payload: string(jsonPayload),
		}
	}

	var data R
	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		return *new(R), err
	}

	return data, nil
}

func PutRequest[T any, R any](url string, payload T) (R, error) {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return *new(R), err
	}

	request, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return *new(R), err
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return *new(R), err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return *new(R), &RequestError{
			Code:    response.StatusCode,
			Message: fmt.Sprintf("HTTP/PUT %s", url),
			Payload: string(jsonPayload),
		}
	}

	var data R
	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		return *new(R), err
	}

	return data, nil
}
