package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sort"
)

// json struct
type SortedNumbers struct {
	Even []int `json:"even"`
	Odd  []int `json:"odd"`
}

func main() {
	http.HandleFunc("/", OddEvenHandler)

	log.Println("Listening on http://localhost:8080 ...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func OddEvenHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is PUT
	if r.Method != http.MethodPut {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error in reading body", http.StatusBadRequest)
		return
	}

	// Decode the JSON body into a slice of integers
	var numbers []int
	if err := json.Unmarshal(body, &numbers); err != nil {
		http.Error(w, "Error in decoding JSON", http.StatusBadRequest)
		return
	}

	even, odd := SplitEvenOdd(numbers)

	// put the calculation into SortedNumbers struct
	result := SortedNumbers{
		Even: even,
		Odd:  odd,
	}

	response, err := json.Marshal(result)
	if err != nil {
		http.Error(w, "Error in creating JSON response", http.StatusInternalServerError)
		return
	}

	// Set header and write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)

}

func SplitEvenOdd(numbers []int) (even, odd []int) {
	// Initialize odd and even slices
	for _, number := range numbers {
		if number%2 == 0 {
			even = append(even, number)
		} else {
			odd = append(odd, number)
		}
	}
	sort.Ints(even)
	sort.Ints(odd)
	return
}
