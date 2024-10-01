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
	gameCommands []Event
	sitCommands []Event
	state *state
	roomName string
}

func createEngine(conn *websocket.Conn, roomName string, smallBlind float64, bigBlind float64) *engine {
	return &engine{
		conn: conn,
		gameCommands: make([]Event, 0),
		sitCommands: make([]Event, 0),
		state: createState(smallBlind, bigBlind, 60),
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

func (e *engine) queueEvent(event Event) {
	if event.EngineCommand == "fold" || event.EngineCommand == "check"  || event.EngineCommand == "call" || event.EngineCommand == "bet" {
		e.gameCommands = append(e.gameCommands, event)
	} else {
		e.sitCommands = append(e.sitCommands, event)
	}
}

func (e *engine) processSitCommand() {
	// copy e.commands so it doesn't change while we're iterating
	commandsCopy := e.sitCommands
	e.sitCommands = make([]Event, 0)
	for _, command := range commandsCopy {
		fmt.Println("processing sit command: ", command)
		seatId := command.SeatId

		if command.EngineCommand == "join" {
			p := createPlayer(command)
			e.state.addPlayer(p)
		} else if command.EngineCommand == "leave" {
			e.state.removePlayer(e.state.players[seatId])
		} else if command.EngineCommand == "startGame" {
			e.startHand()
		} else {
			e.state.players[seatId].makeAction(command)
		}
	}
}

func(e *engine) processGameCommand() {
	// copy e.commands so it doesn't change while we're iterating
	commandsCopy := e.gameCommands
	e.gameCommands = make([]Event, 0)
	for _, command := range commandsCopy {
		fmt.Println("processing game command: ", command)
	}
}

func(e *engine) startHand() {
	e.state.betweenHands = false
	e.postBlinds()
	e.state.spotlight = e.state.dealer.nextPlayerInHand.nextPlayerInHand
}

func(e *engine) postBlinds() {
	smallBlindPlayer := e.state.dealer.nextPlayerInHand
	smallBlindPlayer.postBlind(e.state.smallBlind)
	bigBlindPlayer := e.state.dealer.nextPlayerInHand.nextPlayerInHand
	bigBlindPlayer.postBlind(e.state.bigBlind)
}

func (e engine) sendState() {
	serializeState := createSerializeState(e.state)
	responseMsg, err := json.Marshal(serializeState)
	if err != nil {
		return
	}
	
	e.conn.WriteMessage(websocket.TextMessage, responseMsg)
	fmt.Println("Sending state...")
}