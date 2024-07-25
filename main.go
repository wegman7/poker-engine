package main

import (
	"log"
	"net/http"
)

func main () {
	r := newRoom()

	http.Handle("/ws", r)

	go r.run()

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}