package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type StartGameRequest struct {
	RoomName  string `json:"roomName"`
	SmallBlind  float64 `json:"smallBlind"`
	BigBlind  float64 `json:"bigBlind"`
}

type StartGameResponse struct {
	Message string `json:"message"`
}

func startEngineHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("startEngineHandler")
	req := StartGameRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	
	responseData := StartGameResponse{
		Message: fmt.Sprintf("Started egnine for room %s", req.RoomName),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseData)

	go createEngineConn(req.RoomName, req.BigBlind, req.BigBlind)
}

func main() {
	http.HandleFunc("/start-engine", startEngineHandler)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}