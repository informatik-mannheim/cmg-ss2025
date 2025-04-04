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
	if request.Method == http.MethodPut {
		var numbers []int
		decoder := json.NewDecoder(request.Body)
		// Request isn't successful if the request cannot be decoded into the numbers array
		if err := decoder.Decode(&numbers); err != nil {
			http.Error(responseWriter, "Invalid request body", http.StatusBadRequest)
			return
		}

		evenSlice, oddSlice := separateAndSort(numbers)

		// We build our response body as text
		response := "Even: " + fmt.Sprint(evenSlice) + ", Odd: " + fmt.Sprint(oddSlice)

		// We need to ensure that our response is a plain text
		responseWriter.Header().Set("Content-Type", "application/plain")
		json.NewEncoder(responseWriter).Encode(response)
	} else {
		// We don't handle other request methods
		http.Error(responseWriter, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
}

// This function takes an array of numbers and returns 2 slices containing either even or odd numbers.
// Returns both slices after sorting them.
func separateAndSort(numbers []int) ([]int, []int) {
	evenSlice, oddSlice := []int{}, []int{}
	for number := range numbers {
		if numbers[number]%2 == 0 {
			evenSlice = append(evenSlice, numbers[number])
		} else {
			oddSlice = append(oddSlice, numbers[number])
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
