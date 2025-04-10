package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/", evenOddHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func evenOddHandler(w http.ResponseWriter, r *http.Request) {
	var responseText []byte //Intialize the response

	if r.Method == "PUT"{
		 if r.Header.Get("Content-Type") == "application/json"{
		b, err := io.ReadAll(r.Body) //Reads content of the body and stores it in variable b (byte slice)
		if err != nil {
			http.Error(w, "Error reading Body", 500) //Internal server error
			return
		}

		var m []int
		err = json.Unmarshal(b, &m) //Unmarshals content from byte slice b into int slice m
		if err != nil {
			http.Error(w, "Not an Array of Integers", 400) //Bad Request from Client side
			return
		}

		nums := map[string][]int{ //Map that holds strings as key -> int as value, and contains sorted slices
			"Even": {}, //empty slice
			"Odd":  {},
		}
		
		for _, num := range m {
			if num%2 == 0 { //even-odd check with modolo
				nums["Even"] = append(nums["Even"], num)
			} else {
				nums["Odd"] = append(nums["Odd"], num)
			}
		}
		responseText, _ = json.Marshal(nums) //Marshals content from map -> byte slice
	} else {http.Error(w, "Not in JSON Format", 405)//Status method not allowed
	}//Status method not allowed
	} else {
		http.Error(w, "Not a PUT Request", 405) //Status method not allowed
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Write(responseText)
}
