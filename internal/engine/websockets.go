package engine

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/wegman7/game-engine/config"
)

type Event struct {
    EngineCommand string  `json:"engineCommand"`
    SeatId        int     `json:"seatId"`
	User          string  `json:"user"`
	Chips         float64 `json:"chips"`
}

func deserializeMessage(message []byte) (Event, error) {
	event := Event{}
	err := json.Unmarshal(message, &event)
	if err != nil {
		return event, err
	}
	return event, nil
}

func dial(roomName string, token string) (*websocket.Conn, error) {
	url := fmt.Sprintf("%s/ws/engineconsumer/%s?token=%s", config.AppConfig.BACKEND_URL, roomName, token)
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	return conn, err
}

// readLoop reads messages until the connection closes or a stopEngine command arrives.
// Returns true if stopped cleanly, false on unexpected disconnect.
func readLoop(conn *websocket.Conn, e *engine) bool {
	defer conn.Close()
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("ReadMessage error:", err)
			return false
		}
		event, err := deserializeMessage(message)
		if err != nil {
			fmt.Println("Error decoding JSON:", err)
			return true
		}
		if event.EngineCommand == "stopEngine" {
			return true
		}
		e.queueEvent(event)
	}
}

func CreateEngineConn(roomName string, smallBlind float64, bigBlind float64) {
	token, err := getUserToken(os.Getenv("EMAIL"), os.Getenv("PASSWORD"))
	if err != nil {
		log.Fatal("could not retreive user token:", err)
	}

	const maxRetries = 5
	var e *engine
	stopEngine := make(chan struct{})

	for attempt := range maxRetries {
		if attempt > 0 {
			backoff := time.Duration(1<<attempt) * time.Second
			log.Printf("Reconnecting in %v (attempt %d/%d)...", backoff, attempt+1, maxRetries)
			time.Sleep(backoff)
		}

		conn, err := dial(roomName, token)
		if err != nil {
			log.Printf("Dial failed: %v", err)
			continue
		}

		if e == nil {
			e = createEngine(conn, roomName, smallBlind, bigBlind)
			go e.run(stopEngine)
		} else {
			e.conn = conn
		}

		if clean := readLoop(conn, e); clean {
			close(stopEngine)
			return
		}
	}

	log.Printf("Failed to maintain connection after %d attempts, stopping engine", maxRetries)
	close(stopEngine)
}