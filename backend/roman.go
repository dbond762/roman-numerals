package main

import (
	"net/http"
	"log"
)

func convert(w http.ResponseWriter, r *http.Request) {
	
}


func main() {
	http.HandleFunc("/convert", convert)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
