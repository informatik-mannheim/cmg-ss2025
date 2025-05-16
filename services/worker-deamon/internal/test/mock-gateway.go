package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Hilfsfunktion, um den Request-Body auszugeben
func logRequestBody(r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error reading body:", err)
		return
	}
	defer r.Body.Close()
	if len(body) == 0 {
		fmt.Println("  -> No body received.")
	} else {
		fmt.Println("  -> Payload:", string(body))
	}
}

// Mock Gateway for testing workerdeamon
func main() {
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Received registration")
		logRequestBody(r)
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/worker/heartbeat", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Received heartbeat")
		logRequestBody(r)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"jobs": []map[string]any{
				{
					"id":           "job123",
					"status":       "scheduled",
					"result":       "",
					"errorMessage": "",
				},
				{
					"id":           "job456",
					"status":       "scheduled",
					"result":       "",
					"errorMessage": "",
				},
				{
					"id":           "job789",
					"status":       "scheduled",
					"result":       "",
					"errorMessage": "",
				},
			},
		})
	})

	http.HandleFunc("/result", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Received job result")
		logRequestBody(r)
		w.WriteHeader(http.StatusOK)
	})

	fmt.Println("Mock Gateway listening on :8080")
	http.ListenAndServe(":8080", nil)
}
