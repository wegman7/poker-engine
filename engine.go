package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

var SLEEP_TIME = 5000 * time.Millisecond

type engine struct {
	conn *websocket.Conn
	commands *[][]byte
	state *state
}

func createEngine(conn *websocket.Conn) *engine {
	commands := make([][]byte, 0)
	return &engine{
		conn: conn,
		commands: &commands,
		state: createState(),
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

func (e *engine) tick() {
	// copy e.commands so it doesn't change while we're iterating
	commandsCopy := *e.commands
	*e.commands = make([][]byte, 0)
	for _, command := range commandsCopy {
		e.processCommand(command)
	}
}

func (e *engine) queueCommand(command []byte) {
	*e.commands = append(*e.commands, command)
}

func(e *engine) processCommand(commandBytes []byte) {
	command := string(commandBytes)
	fmt.Println(command)
}

func (e engine) sendState() {
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