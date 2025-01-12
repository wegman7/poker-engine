package engine

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type Event struct {
    EngineCommand string  `json:"engineCommand"`
    SeatId        int     `json:"seatId"`
	User          string  `json:"user"`
	Chips         float64 `json:"chips"`
}

func closeConn(conn *websocket.Conn, stopEngine chan struct{}, roomName string) {
	// stop engine goroutine
	close(stopEngine)
	log.Println("Closing websockets connection for room", roomName)
	conn.Close()
}

func deserializeMessage(message []byte) (Event, error) {
	event := Event{}
	err := json.Unmarshal(message, &event)
	if err != nil {
		return event, err
	}
	return event, nil
}

func CreateEngineConn(roomName string, smallBlind float64, bigBlind float64) {
	url := fmt.Sprintf("ws://localhost:8000/ws/engineconsumer/%s/", roomName)

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	e := createEngine(conn, roomName, smallBlind, bigBlind)
	stopEngine := make(chan struct{})
	go e.run(stopEngine)
    defer closeConn(conn, stopEngine, roomName)

    for {
        _, message, err := conn.ReadMessage()
        if err != nil {
            log.Println("ReadMessage error:", err)
            break
        }

		event, err := deserializeMessage(message)
		if err != nil {
			fmt.Println("Error decoding JSON:", err)
			return
		}
		if event.EngineCommand == "stopEngine" {
			break
		}

		e.queueEvent(event)
    }
}