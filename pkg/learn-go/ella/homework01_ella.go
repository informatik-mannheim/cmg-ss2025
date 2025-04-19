package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"slices"
)

func main() {
	// Register HTTP handler and start server on port 8080
	http.HandleFunc("/", evenOddHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// evenOddHandler handles incoming HTTP requests
func evenOddHandler(w http.ResponseWriter, r *http.Request) {
	var responseText []byte // Will hold the response to be sent back

	// Ensure the request uses PUT method
	if r.Method == "PUT" {

		// Check that the content type is application/json
		if r.Header.Get("Content-Type") == "application/json" {

			// Read the entire request body
			b, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Error reading Body", 500) // Internal server error
				return
			}

			// Declare slice to hold incoming numbers
			var m []int

			// Attempt to parse JSON body into integer slice
			err = json.Unmarshal(b, &m)
			if err != nil {
				http.Error(w, "Not an Array of Integers", http.StatusBadRequest) // Bad client input
				return
			}

			// Process the numbers into even and odd groups
			nums := sortNums(m)

			// Marshal the map into JSON to return as response
			responseText, _ = json.Marshal(nums)

		} else {
			// If the content type is not JSON
			http.Error(w, "Not in JSON Format", 405)
		}

	} else {
		// If the method is not PUT
		http.Error(w, "Not a PUT Request", 405)
	}

	// Set response type to plain text (even though we're sending JSON)
	w.Header().Set("Content-Type", "text/plain")
	w.Write(responseText) // Write response back to client
}

// sortNums separates even and odd integers and sorts them
func sortNums(m []int) map[string][]int {

	// Create a map with two slices to categorize numbers
	nums := map[string][]int{
		"Even": {}, // Holds even numbers
		"Odd":  {}, // Holds odd numbers
	}

	// Iterate over all numbers and categorize them
	for _, num := range m {
		if num%2 == 0 {
			nums["Even"] = append(nums["Even"], num)
		} else {
			nums["Odd"] = append(nums["Odd"], num)
		}
	}
	// Sort each list after insertion
	slices.Sort(nums["Even"])
	slices.Sort(nums["Odd"])
	// Return the categorized and sorted map
	return nums
}
