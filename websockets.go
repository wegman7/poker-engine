package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
)

const (
	socketBufferSize  = 1024
	messageBufferSize = 1024
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

func parseRequestParams(req *http.Request) (string, float64, error) {
	// Extract roomName from the path
	roomName := strings.TrimPrefix(req.URL.Path, "/ws/")

	// Extract bigBlind from the query parameters and convert to float64
	query := req.URL.Query()
	bigBlind := query.Get("bigBlind")
	bigBlindFloat, err := strconv.ParseFloat(bigBlind, 64)
	if err != nil {
		return "", 0, fmt.Errorf("failed to convert bigBlind to float64: %w", err)
	}

	return roomName, bigBlindFloat, nil
}

func ServeHTTP(w http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("Serve http: ", err)
		return
	}

	roomName, bigBlindFloat, err := parseRequestParams(req)
	if err != nil {
		log.Fatal("Failed to parse request params: ", err)
		return
	}
	e := createEngine(conn, roomName, bigBlindFloat)

	stopEngine := make(chan struct{})
	go e.run(stopEngine)

    defer func() {
		// stop engine goroutine
		close(stopEngine)
		fmt.Println("stopping engine")
        conn.Close()
    }()
    for {
        _, message, err := conn.ReadMessage()
        if err != nil {
            log.Println("ReadMessage:", err)
            break
        }

		// Declare a map to hold the decoded data
		var messageMap map[string]string

		// Unmarshal (convert) the JSON into the map
		err2 := json.Unmarshal(message, &messageMap)
		if err2 != nil {
			fmt.Println("Error decoding JSON:", err2)
			return
		}

		fmt.Println("new message: ", messageMap)
		e.queueCommand(messageMap)
    }
}