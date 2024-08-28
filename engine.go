package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

var SLEEP_TIME = 200 * time.Millisecond

type engine struct {
	conn *websocket.Conn
	gameCommands [][]byte
	state state
}

func createEngine(conn *websocket.Conn) *engine {
	return &engine{
		conn: conn,
		gameCommands: make([][]byte, 0),
		state: *createState(),
	}
}

func (e *engine) run(stopEngine chan struct{}) {
	for {
		select {
		case <-stopEngine:
			log.Println("Goroutine stopped as WebSocket is closed.")
            return
		default:
			time.Sleep(SLEEP_TIME)
			e.tick()
			e.sendState()
		}
	}
}

func (e engine) tick() {
	return
}

func (e engine) queueGameCommand() {
	for {
		_, msg, err := e.conn.ReadMessage()
		if err != nil {
			return
		}
		e.gameCommands = append(e.gameCommands, msg)
		fmt.Println("adding command to queue", string(e.gameCommands[0]))
	}
}

func (e engine) sendState() {
	fmt.Println("sending state...", e)
	serializePlayer := SerializePlayer{
		SeatId: "1",
	}
	serializeState := SerializeState{
		Type: "gamestate",
		Players: map[string]SerializePlayer{
			"player1": serializePlayer,
		},
	}

	responseMsg, err := json.Marshal(serializeState)
	if err != nil {
		return
	}
	e.conn.WriteMessage(websocket.TextMessage, responseMsg)
}