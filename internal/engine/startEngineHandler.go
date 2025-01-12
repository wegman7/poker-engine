package engine

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

func StartEngineHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("startEngineHandler")
	req := StartGameRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	
	responseData := StartGameResponse{
		Message: fmt.Sprintf("Started engine for room %s", req.RoomName),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseData)

	go CreateEngineConn(req.RoomName, req.SmallBlind, req.BigBlind)
}