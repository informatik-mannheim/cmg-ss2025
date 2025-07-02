package gateway

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"worker-daemon/internal/ports"
)

type Client struct {
	BaseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL:    baseURL,
		httpClient: &http.Client{},
	}
}

func checkStatusOK(resp *http.Response) error {
	if resp.StatusCode != http.StatusOK {
		return errors.New("http error: status code " + http.StatusText(resp.StatusCode))
	}
	return nil
}

func (c *Client) Register(key string, zone string) (*ports.RegisterResponse, error) {
	payload := map[string]string{
		"key":  key,
		"zone": zone,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(c.BaseURL+"/register", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := checkStatusOK(resp); err != nil {
		return nil, err
	}
	var regResp *ports.RegisterResponse
	if err := json.NewDecoder(resp.Body).Decode(&regResp); err != nil {
		return nil, err
	}

	return regResp, err
}

func (c *Client) SendHeartbeat(workerId string, status string, token string) ([]ports.Job, error) {
	payload := map[string]string{
		"workerId": workerId,
		"status":   status,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.BaseURL+"/worker/heartbeat", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := checkStatusOK(resp); err != nil {
		return nil, err
	}

	var jobs []ports.Job
	if err := json.NewDecoder(resp.Body).Decode(&jobs); err != nil {
		return nil, err
	}

	return jobs, nil
}

func (c *Client) SendResult(j ports.Job, token string) error {
	payload := map[string]string{
		"jobId":        j.ID,
		"status":       j.Status,
		"result":       j.Result,
		"errorMessage": j.ErrorMessage,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", c.BaseURL+"/result", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return checkStatusOK(resp)
}
