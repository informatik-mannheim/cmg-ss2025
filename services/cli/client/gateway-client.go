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
	resp, err := http.Post("http://localhost:"+port+"/jobs", "application/json", bytes.NewBuffer(jsonRequest))
	if err != nil {
		log.Fatal("Error creating job", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
}

func GetJobById(id string) {
	url := fmt.Sprintf("http://localhost:%s/jobs/%s", port, id)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("Error getting job:", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

}

func GetJobOutcome(id string) {

	url := fmt.Sprintf("http://localhost:%s/jobs/%s/outcome", port, id)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("Error getting job outcome:", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

}
