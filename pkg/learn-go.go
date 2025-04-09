package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)




func main(){

http.HandleFunc("/", evenOddHandler)
log.Fatal(http.ListenAndServe(":8080", nil))
}

func evenOddHandler(w http.ResponseWriter, r *http.Request){

if r.Method == "PUT"{
	b, err := io.ReadAll(r.Body) //Reads content of the body and stores it in variable b (byte slice)
	if err != nil { 
		http.Error(w, "Error reading Body", 500) //Internal server error
		return
	} 

	var m []int
	err = json.Unmarshal(b, &m) //unmarshals content from byte slice b into int slice m
	if err != nil {
		http.Error(w, "Error converting JSON", 400) //Bad Request from Client side
		return
	}

	nums := map[string][]int{ //Map that holds strings as key -> int as value, and contains sorted slices
		"Even": {},//empty slice
		"Odd": {}, 
	}

	for i:= 0; i < len(m); i++ {
		if i%2 == 0 { //even-odd check with modolo
			nums["Even"] = append(nums["Even"], i)
		} else {nums["Odd"] = append(nums["Odd"], i)}
	}

	} else { http.Error(w, "Not a PUT Request", 405)}//Status method not allowed

	w.Header().Set("Content-Type", "text")
	w.WriteHeader(200)// Status OK

	json.NewEncoder(w).Encode{
	}

}
