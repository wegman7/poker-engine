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

var runningEngines = make(map[string]struct{})

func StartEngineHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("startEngineHandler")
	req := StartGameRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if _, ok := runningEngines[req.RoomName]; ok {
		http.Error(w, "Engine already running for room", http.StatusBadRequest)
		return
	}
	runningEngines[req.RoomName] = struct{}{}
	go CreateEngineConn(req.RoomName, req.SmallBlind, req.BigBlind)
	
	responseData := StartGameResponse{
		Message: fmt.Sprintf("Started engine for room %s", req.RoomName),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseData)
}