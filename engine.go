package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

var SLEEP_TIME = 100 * time.Millisecond

type engine struct {
	conn *websocket.Conn
	gameCommands []map[string]string
	sitCommands []map[string]string
	state *state
	roomName string
}

func createEngine(conn *websocket.Conn, roomName string, bigBlind float64) *engine {
	return &engine{
		conn: conn,
		gameCommands: make([]map[string]string, 0),
		sitCommands: make([]map[string]string, 0),
		state: createState(bigBlind, 60),
		roomName: roomName,
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
	if e.state.betweenHands {
		e.processSitCommand()
	} else {
		e.processGameCommand()
	}
}

func (e *engine) queueCommand(command map[string]string) {
	if command["engineCommand"] == "fold" || command["engineCommand"] == "check"  || command["engineCommand"] == "call" || command["engineCommand"] == "bet" {
		e.gameCommands = append(e.gameCommands, command)
	} else {
		e.sitCommands = append(e.sitCommands, command)
	}
}

func (e *engine) processSitCommand() {
	// copy e.commands so it doesn't change while we're iterating
	commandsCopy := e.sitCommands
	e.sitCommands = make([]map[string]string, 0)
	for _, command := range commandsCopy {
		fmt.Println("processing sit command: ", command)
	}
}

func(e *engine) processGameCommand() {
	// copy e.commands so it doesn't change while we're iterating
	commandsCopy := e.gameCommands
	e.gameCommands = make([]map[string]string, 0)
	for _, command := range commandsCopy {
		fmt.Println("processing game command: ", command)
	}
}

func (e engine) sendState() {
	serializePlayer := SerializePlayer{
		SeatId: "1",
	}
	serializeState := SerializeState{
		ChannelCommand: "sendState",
		RoomName: e.roomName,
		Players: map[string]SerializePlayer{
			"player1": serializePlayer,
		},
	}

	responseMsg, err := json.Marshal(serializeState)
	if err != nil {
		return
	}
	e.conn.WriteMessage(websocket.TextMessage, responseMsg)
	fmt.Println("Sending state...")
}