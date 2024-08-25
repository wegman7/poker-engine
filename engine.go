package main

import (
	"fmt"
	"time"
	"github.com/gorilla/websocket"
	"encoding/json"
)

var SLEEP_TIME = 200 * time.Millisecond

type engine struct {
	socket *websocket.Conn
	gameCommands [][]byte
	state state
	running bool
}

func startEngine(socket *websocket.Conn) *engine {
	e := &engine{
		socket: socket,
		gameCommands: make([][]byte, 0),
		state: *createState(),
		running: true,
	}

	i := 0
	for i < 10 {
		fmt.Println("sending hello")
		time.Sleep(SLEEP_TIME)
		e.socket.WriteMessage(websocket.TextMessage, []byte("hello"))
		i++
	}
	// go e.run()
	// go e.queueGameCommand()

	return e
}

func (e *engine) run() {
	for e.running {
		time.Sleep(SLEEP_TIME)
		e.tick()
		e.sendState()
	}
	fmt.Println("Engine stopped")
}

func (e engine) tick() {
	return
}

func (e engine) queueGameCommand() {
	for {
		_, msg, err := e.socket.ReadMessage()
		if err != nil {
			return
		}
		e.gameCommands = append(e.gameCommands, msg)
		fmt.Println("adding command to queue", string(e.gameCommands[0]))
	}
}

func (e engine) sendState() {
	fmt.Println("sending state...", e.running, e)
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
	e.socket.WriteMessage(websocket.TextMessage, responseMsg)
}