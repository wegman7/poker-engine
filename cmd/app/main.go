package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/wegman7/game-engine/config"
	"github.com/wegman7/game-engine/internal/engine"
)

func main() {
	// Use a command-line flag or environment variable to determine the environment
	env := flag.String("env", "dev", "Environment to run: dev or prod")
	flag.Parse()

	// Load the configuration based on the chosen environment
	if err := config.Load(*env); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	http.HandleFunc("/start-engine", engine.StartEngineHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}