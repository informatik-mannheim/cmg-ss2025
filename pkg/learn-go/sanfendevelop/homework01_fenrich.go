package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
)

type RequestData struct {
	Numbers []int
}

func main() {
	http.HandleFunc("/numbers/sorted", handler)
	fmt.Println("Server is listening on port 8080")
	// make sure that localhost is used
	http.ListenAndServe("127.0.0.1:8080", nil)
}

// Takes the content from the request body and saves it in a structure with a int-slice

func decodeRequest(r *http.Request) (RequestData, error) {
	var data RequestData
	// if the contents cannot be assigned, error is returned
	// e. g. float input 1.1, 1.2
	if err := json.NewDecoder(r.Body).Decode(&data.Numbers); err != nil {
		return RequestData{}, err
	}
	return data, nil
}

// Sorts numbers according to whether they are even or odd

func sortNumbers(numbers []int) (oddNumbers, evenNumbers []int) {
	for _, num := range numbers {
		if num%2 == 0 {
			evenNumbers = append(evenNumbers, num)
		} else {
			oddNumbers = append(oddNumbers, num)
		}
	}

	sort.Ints(oddNumbers)
	sort.Ints(evenNumbers)
	return // function signature determines which values are returned
}

// Converts number sequences into character strings

func encodeResponse(oddNumbers, evenNumbers []int) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Odd numbers: %v\n", oddNumbers))
	sb.WriteString(fmt.Sprintf("Even numbers: %v\n", evenNumbers))
	return sb.String()
}

// Handles incoming request as Json and creates response as text

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	data, err := decodeRequest(r)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	oddNumbers, evenNumbers := sortNumbers(data.Numbers)
	response := encodeResponse(oddNumbers, evenNumbers)

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(response))
}
