package gateway

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type Client struct {
	BaseURL string
}

func NewClient(baseURL string) *Client {
	return &Client{BaseURL: baseURL}
}

func (c *Client) Register(workerId, key, location string) error {
	payload := map[string]string{
		"id":       workerId,
		"key":      key,
		"location": location,
	}
	return c.postJSON("/register", payload)
}

func (c *Client) SendHeartbeat(workerId, status string) error {
	payload := map[string]string{
		"workerId": workerId,
		"status":   status,
	}
	return c.postJSON("/worker/heartbeat", payload)
}

func (c *Client) ReportResult(jobId, status, result, errorMsg string) error {
	payload := map[string]string{
		"jobId":        jobId,
		"status":       status,
		"result":       result,
		"errorMessage": errorMsg,
	}
	return c.postJSON("/result", payload)
}

func (c *Client) postJSON(path string, payload any) error {
	data, _ := json.Marshal(payload)
	resp, err := http.Post(c.BaseURL+path, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
