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

b, err := io.ReadAll(r.Body)
if err != nil {
	http.Error(w, "Error reading Body", 500)
} 

var m map[string]int
c, err := json.Unmarshal(b, &m) 


} else { http.Error(w, "Not a PUT Request", 500)}
}
