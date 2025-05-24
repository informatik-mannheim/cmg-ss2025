package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

var port string

func init() {
	port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
}

func CreateJob(image_id string, job_name string, creation_zone string, parameters map[string]string) {
	request := CreateJobRequest{
		ImageID:      image_id,
		JobName:      job_name,
		CreationZone: creation_zone,
		Parameters:   parameters}

	jsonRequest, err := json.Marshal(request)
	if err != nil {
		log.Fatal("Error creating job", err)
	}
	resp, err := http.Post("http://localhost:"+port+"/jobs", "application/json", bytes.NewBuffer(jsonRequest))
	if err != nil {
		log.Fatal("Error creating job", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	log.Println("Response status: ", resp.Status)
	log.Println("Response body: ", string(body))
}

func GetJobById(id string) {
	url := fmt.Sprintf("http://localhost:%s/jobs/%s", port, id)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("Error getting job:", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Status:", resp.Status)
	fmt.Println("Response:", string(body))
}

func GetJobOutcome(id string) {

	url := fmt.Sprintf("http://localhost:%s/jobs/%s/outcome", port, id)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("Error getting job outcome:", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Status:", resp.Status)
	fmt.Println("Response:", string(body))
}
