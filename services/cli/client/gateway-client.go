package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/informatik-mannheim/cmg-ss2025.git"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

var port string
var AuthToken string

type GatewayClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewGatewayClient(baseURL string) *GatewayClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &GatewayClient{baseURL: baseURL, httpClient: &http.Client{}}
}

func init() {
	port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
}

func (c *GatewayClient) CreateJob(jobName string, creationZone string, imageId cli.ContainerImage, parameters map[string]string) {
	request := cli.CreateJobRequest{
		JobName:      jobName,
		CreationZone: creationZone,
		Image:        imageId,
		Parameters:   parameters}

	jsonRequest, err := json.Marshal(request)
	if err != nil {
		log.Fatal("Error creating job", err)
	}

	url := c.baseURL + "/jobs"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonRequest))
	if err != nil {
		log.Fatal("Error creating request", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+AuthToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error creating job", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(resp.Body)

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Status:", resp.Status)
	fmt.Println("Response:", string(body))

}

func (c *GatewayClient) GetJobById(id string) {
	url := fmt.Sprintf("%s/jobs/%s", c.baseURL, id)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("Error creating request:", err)
	}

	req.Header.Set("Authorization", "Bearer "+AuthToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error making request:", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Status:", resp.Status)
	fmt.Println("Response:", string(body))
}

func (c *GatewayClient) GetJobOutcome(id string) {
	url := fmt.Sprintf("%s/jobs/%s/outcome", c.baseURL, id)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("Error creating request:", err)
	}

	req.Header.Set("Authorization", "Bearer "+AuthToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error making request:", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Status:", resp.Status)
	fmt.Println("Response:", string(body))

}

func (c *GatewayClient) Login(secret string) {
	url := fmt.Sprintf("%s/auth/login", c.baseURL)

	payload := map[string]string{
		"secret": secret,
	}
	jsonBody, err := json.Marshal(payload)
	if err != nil {
		log.Fatal("Error marshalling JSON:", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Fatal("Error making POST request:", err)
	}
	defer resp.Body.Close()

	var tokenResp cli.TokenResponse

	AuthToken = tokenResp.Token
	fmt.Println("Status:", resp.Status)
	if resp.StatusCode == 200 {
		fmt.Println("Login successful.")
	} else {
		fmt.Println("Unsuccessful Login.")
	}

}
