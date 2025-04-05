package tonka

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"syscall"
)

type SortedNumbers struct {
	Even []int `json:"even"`
	Odd  []int `json:"odd"`
}

func main() {
	server := &http.Server{Addr: ":8080"}

	http.HandleFunc("/", Handler)

	// graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		log.Print("The service is shutting down...")
		server.Shutdown(context.Background())
	}()

	log.Print("listening...")
	server.ListenAndServe()
	log.Print("Done")
}

func Handler(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPut {
		http.Error(response, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var numbers []int
	decoder := json.NewDecoder(request.Body)
	if err := decoder.Decode(&numbers); err != nil {
		http.Error(response, "Invalid request body", http.StatusBadRequest)
		return
	}
	sortedNumbers := SeperateAndSort(numbers)

	response.Header().Set("Content-Type", "application/json")
	json.NewEncoder(response).Encode(sortedNumbers)
	response.WriteHeader(http.StatusOK)
	log.Print("Request processed successfully")
}

func SeperateAndSort(numbers []int) SortedNumbers {
	evenSlice, oddSlice := []int{}, []int{}
	sortedInput := numbers

	sort.Ints(sortedInput)

	for _, number := range numbers {
		if number%2 == 0 {
			evenSlice = append(evenSlice, number)
		} else {
			oddSlice = append(oddSlice, number)
		}
	}

	sortedNumbers := SortedNumbers{
		Even: evenSlice,
		Odd:  oddSlice,
	}

	return sortedNumbers
}
