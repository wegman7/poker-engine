package engine

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gorilla/websocket"
	"github.com/wegman7/game-engine/config"
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
	token, err := getUserToken(os.Getenv("EMAIL"), os.Getenv("PASSWORD"))
	if err != nil {
		log.Fatal("could not retreive user token:", err)
	}
	url := fmt.Sprintf("ws://%s/ws/engineconsumer/%s?token=%s", roomName, config.AppConfig.BACKEND_URL, token)

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