package main

import (
	"net/http"
	"log"
)

func Convert(w http.ResponseWriter, r *http.Request) {
	
}


func main() {
	http.HandleFunc("/convert", Convert)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
