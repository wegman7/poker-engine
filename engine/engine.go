package engine

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/chehsunliu/poker"
	"github.com/gorilla/websocket"
)

var SLEEP_TIME = 100 * time.Millisecond

type engineState int

const (
	StateProcessSitCommands engineState = iota
	StateProcessGameCommands
	StateStartHand
	StateEndStreet
	StateEveryoneFolded
)

type engine struct {
	conn         *websocket.Conn
	gameCommands []Event
	sitCommands  []Event
	state        *state
	roomName     string
}

func createEngine(conn *websocket.Conn, roomName string, smallBlind float64, bigBlind float64) *engine {
	return &engine{
		conn:         conn,
		gameCommands: make([]Event, 0),
		sitCommands:  make([]Event, 0),
		state:        createState(smallBlind, bigBlind, 60),
		roomName:     roomName,
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
	// use states here
	if e.state.handInAction {
		e.processGameCommand()
		if e.isStreetComplete() {
			e.endStreet()
		}
	} else {
		e.processSitCommand()
	}
}

func (e *engine) queueEvent(event Event) {
	if event.EngineCommand == "fold" || event.EngineCommand == "check" || event.EngineCommand == "call" || event.EngineCommand == "bet" {
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
		user := command.User

		if command.EngineCommand == "join" {
			p := createPlayer(command)
			e.state.addPlayer(p)
		} else if command.EngineCommand == "leave" {
			e.state.removePlayer(e.state.players[user])
		} else if command.EngineCommand == "startGame" {
			e.startHand()
		} else {
			e.state.players[user].makeAction(command, e, e.state)
		}
	}
}

func (e *engine) processGameCommand() {
	// copy e.commands so it doesn't change while we're iterating
	commandsCopy := e.gameCommands
	e.gameCommands = make([]Event, 0)
	for _, command := range commandsCopy {
		fmt.Println("processing game command: ", command)
	}
}

func (e *engine) startHand() {
	if err := e.state.performDealerRotation(); err != nil {
		log.Println("Error rotating dealer: ", err)
		return
	}
	e.postBlinds()
	e.dealCards()
	e.state.handInAction = true
}

func (e *engine) postBlinds() {
	playerCount := e.state.countPlayersInHand()
	// when a hand is heads up, the dealer posts the sb and goes first preflop
	if playerCount == 2 {
		e.state.dealer.postBlind(e.state.smallBlind)
		e.state.dealer.nextInHand.postBlind(e.state.bigBlind)
		e.state.spotlight = e.state.dealer
		return
	}
	e.state.dealer.nextInHand.postBlind(e.state.smallBlind)
	e.state.dealer.nextInHand.nextInHand.postBlind(e.state.bigBlind)
	e.state.spotlight = e.state.dealer.nextInHand.nextInHand.nextInHand
}

func (e *engine) dealCards() {
	deck := poker.NewDeck()
	pointer := e.state.dealer.nextInHand
	for {
		pointer.holeCards = deck.Draw(2)
		if pointer == e.state.dealer {
			break
		}
		pointer = pointer.nextInHand
	}
}

func (e *engine) isStreetComplete() bool {
	return e.state.spotlight == e.state.lastAggressor
}

func (e *engine) everyoneFolded() {
	winner := e.state.psuedoDealer

	chips := 0.0
	for _, player := range e.state.players {
		chips += player.chipsInPot
		player.chipsInPot = 0
	}
	chips += e.state.pot

	winner.chips += chips
}

func (e *engine) endStreet() {
	// check if we have enough players to continue

	e.state.spotlight = e.state.dealer
	e.state.lastAggressor = nil
	// deal next street, might need to add street as state attribute
}

func (e *engine) dealStreet() {
	// deal next street
}

func (e *engine) showdown() {
	// determine winner
}

func (e *engine) payout() {
	// should call endHand
}

func (e *engine) endHand() {
	// the opposite of startHand
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
