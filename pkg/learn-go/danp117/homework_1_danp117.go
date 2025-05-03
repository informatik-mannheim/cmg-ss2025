package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
)

// main function listening to port 8080
func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// processes PUT request and responds with sorted map as string
func handler(w http.ResponseWriter, r *http.Request) {

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Invalid Content-Type. Expected application/json", http.StatusUnsupportedMediaType)
		return
	}
  
	if r.Method != http.MethodPut {
		http.Error(w, "This method is not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error while reading request body", http.StatusBadRequest)
		return
	}

	var numbers []int
	if err := json.Unmarshal(body, &numbers); err != nil {
		http.Error(w, "Invalid JSON input", http.StatusBadRequest)
		return
	}

	result := SeparateEvenOdd(numbers)

	response := fmt.Sprintf("even : %s, odd : %s", formatSliceWithComma(result["even"]), formatSliceWithComma(result["odd"]))
	w.Header().Set("Content-Type", "text/plain")
	_, _ = w.Write([]byte(response))
}

// takes a slice of integers and returns a map with sorted even and odd numbers.
func SeparateEvenOdd(numbers []int) map[string][]int {
	result := make(map[string][]int)
	for _, num := range numbers {
		if num%2 == 0 {
			result["even"] = append(result["even"], num)
		} else {
			result["odd"] = append(result["odd"], num)
		}
	}
	sort.Ints(result["even"])
	sort.Ints(result["odd"])
	return result
}

func formatSliceWithComma(nums []int) string {
	if len(nums) == 0 {
		return "[]"
	}
	s := "["
	for i, n := range nums {
		if i > 0 {
			s += ", "
		}
		s += fmt.Sprintf("%d", n)
	}
	s += "]"
	return s
}
