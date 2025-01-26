package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/wegman7/game-engine/internal/engine"
)

func main() {
	err := godotenv.Load()
	if err != nil {
	  log.Fatal("Error loading .env file")
	}

	http.HandleFunc("/start-engine", engine.StartEngineHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}