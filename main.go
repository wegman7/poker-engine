package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/ws/", ServeHTTP)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}