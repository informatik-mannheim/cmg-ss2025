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
)

var port string
var AuthToken string

func init() {
	port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
}

func CreateJob(jobName string, creationZone string, imageId cli.ContainerImage, parameters map[string]string) {
	request := cli.CreateJobRequest{
		JobName:      jobName,
		CreationZone: creationZone,
		Image:        imageId,
		Parameters:   parameters}

	jsonRequest, err := json.Marshal(request)
	if err != nil {
		log.Fatal("Error creating job", err)
	}

	url := "http://localhost:" + port + "/jobs"

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

func GetJobById(id string) {
	url := fmt.Sprintf("http://localhost:%s/jobs/%s", port, id)

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

func GetJobOutcome(id string) {
	url := fmt.Sprintf("http://localhost:%s/jobs/%s/outcome", port, id)

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

func Login(secret string) {
	url := fmt.Sprintf("http://localhost:%s/auth/login", port)

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
