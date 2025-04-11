package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
)

func handler(responseWriter http.ResponseWriter, request *http.Request) {
	// We only accept PUT requests
	if request.Method != http.MethodPut {
		http.Error(responseWriter, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	var numbers []int
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&numbers)

	if err != nil {
		http.Error(responseWriter, "Invalid request body", http.StatusBadRequest)
		return
	}

	evenSlice, oddSlice := separateAndSort(numbers)

	// We build our response body as text
	response := "Even: " + fmt.Sprint(evenSlice) + ", Odd: " + fmt.Sprint(oddSlice)

	// We need to ensure that our response is a plain text
	responseWriter.Header().Set("Content-Type", "application/plain")
	responseWriter.WriteHeader(http.StatusOK)
	responseWriter.Write([]byte(response))

}

// This function takes an array of numbers and returns 2 slices containing either even or odd numbers.
// Returns both slices after sorting them.
func separateAndSort(numbers []int) ([]int, []int) {
	evenSlice, oddSlice := []int{}, []int{}
	for _, number := range numbers {
		if number%2 == 0 {
			evenSlice = append(evenSlice, number)
		} else {
			oddSlice = append(oddSlice, number)
		}
	}

	sort.Ints(evenSlice)
	sort.Ints(oddSlice)
	return evenSlice, oddSlice
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Server is listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
