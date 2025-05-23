package client

import (
	"bytes"
	"encoding/json"
	"github.com/informatik-mannheim/cmg-ss2025/services/consumer-gateway/ports"
	"io"
	"log"
	"net/http"
	"os"
)

func CreateJob(image_id string, job_name string, creation_zone string, parameters map[string]string) {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	request := ports.CreateJobRequest{
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
