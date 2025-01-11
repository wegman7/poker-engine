package main

import (
	"log"
	"net/http"
	"github.com/wegman7/game-engine/internal/engine"
)

func main() {
	http.HandleFunc("/start-engine", engine.StartEngineHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}