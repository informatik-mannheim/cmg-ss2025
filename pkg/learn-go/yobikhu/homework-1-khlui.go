package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

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

	odd, even := SplitOddEven(numbers)
	// Sort the odd and even slices
	sort.Ints(even)
	sort.Ints(odd)

	// Check if the slices are empty
	if len(odd) == 0 && len(even) == 0 {
		http.Error(w, "No numbers provided", http.StatusBadRequest)
		return
	}

	response := fmt.Sprintf("even: %v, odd: %v", formatSlice(even), formatSlice(odd))
	w.Write([]byte(response))

}

func SplitOddEven(numbers []int) (odd []int, even []int) {
	// Initialize odd and even slices
	for _, number := range numbers {
		if number%2 == 0 {
			even = append(even, number)
		} else {
			odd = append(odd, number)
		}
	}
	return odd, even
}

// Extended function to make the output like the expected output

// Converts a slice of ints to a formatted string: [1,2,3]
func formatSlice(nums []int) string {
	strs := make([]string, len(nums))
	for i, v := range nums {
		strs[i] = strconv.Itoa(v)
	}
	return "[" + strings.Join(strs, ",") + "]"
}
