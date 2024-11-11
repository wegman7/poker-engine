package engine

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/chehsunliu/poker"
	"github.com/gorilla/websocket"
)

var ENGINE_LOOP_PAUSE = 100 * time.Millisecond
var PAUSE_SHORT = 100 * time.Millisecond
var PAUSE_MEDIUM = 200 * time.Millisecond
var PAUSE_LONG = 300 * time.Millisecond

type engineState int

const (
	StateProcessSitCommands engineState = iota
	StateStartHand
	StatePauseAfterStartHand
	StatePostBlinds
	StatePauseAfterPostBlinds
	StateDealCards
	StateProcessGameCommands
	StatePauseAfterEveryoneFolded
	StateEveryoneFoldedPayout
	StatePauseAfterEveryoneFoldedPayout
	StateEndStreet
	StatePauseAfterEndStreet
	StatePayout
	StatePauseAfterStatePayout
	StateShowdownPayout
	StateEndHand
	StateDealNextStreet
)

type engine struct {
	conn         *websocket.Conn
	gameCommands []Event
	sitCommands  []Event
	state        *state
	roomName     string
	engineState  engineState
}

func createEngine(conn *websocket.Conn, roomName string, smallBlind float64, bigBlind float64) *engine {
	return &engine{
		conn:         conn,
		gameCommands: make([]Event, 0),
		sitCommands:  make([]Event, 0),
		state:        createState(smallBlind, bigBlind, 60),
		roomName:     roomName,
		engineState:  StateProcessSitCommands,
	}
}

func (e *engine) run(stopEngine chan struct{}) {
	e.transitionState(StateProcessSitCommands)
	for {
		select {
		case <-stopEngine:
			log.Println("Goroutine stopped as WebSocket is closed.")
			return
		default:
			time.Sleep(ENGINE_LOOP_PAUSE)
			e.tick()
			e.sendState()
		}
	}
}

func (e *engine) tick() {
	// use states here
	switch e.engineState {
	case StateProcessSitCommands:
		e.processSitCommand()
	case StateProcessGameCommands:
		e.processGameCommand()
	case StateStartHand:
		e.startHand()
	case StatePauseAfterStartHand:
		e.pauseAfterStartHand()
	case StatePostBlinds:
		e.postBlinds()
	case StatePauseAfterPostBlinds:
		e.pauseAfterPostBlinds()
	case StateDealCards:
		e.dealCards()
	case StateEndStreet:
		e.endStreet()
	case StatePauseAfterEveryoneFolded:
		e.pauseAfterEveryoneFolded()
	case StateEveryoneFoldedPayout:
		e.everyoneFoldedPayout()
	case StatePauseAfterEveryoneFoldedPayout:
		e.pauseAfterEveryoneFoldedPayout()
	case StateEndHand:
		e.endHand()
	}
}

func (e *engine) transitionState(newEngineState engineState) {
	log.Println("Transitioning state from", e.engineState, "to", newEngineState)
	e.engineState = newEngineState
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
			e.transitionState(StateStartHand)
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
		user := command.User
		e.state.players[user].makeAction(command, e, e.state)
	}
}

func (e *engine) startHand() {
	if err := e.state.performDealerRotation(); err != nil {
		log.Println("Error rotating dealer: ", err)
		return
	}
	e.transitionState(StatePauseAfterStartHand)
}

func (e *engine) pauseAfterStartHand() {
	time.Sleep(5 * time.Second)
	e.transitionState(StatePostBlinds)
}

func (e *engine) postBlinds() {
	playerCount := e.state.countPlayersSittingIn()
	// when a hand is heads up, the dealer posts the sb and goes first preflop
	if playerCount == 2 {
		e.state.dealer.postBlind(e.state.smallBlind)
		e.state.dealer.nextInHand.postBlind(e.state.bigBlind)
		e.state.spotlight = e.state.dealer
		e.transitionState(StatePauseAfterPostBlinds)
		return
	}
	e.state.dealer.nextInHand.postBlind(e.state.smallBlind)
	e.state.dealer.nextInHand.nextInHand.postBlind(e.state.bigBlind)
	e.state.spotlight = e.state.dealer.nextInHand.nextInHand.nextInHand

	e.transitionState(StatePauseAfterPostBlinds)
}

func (e *engine) pauseAfterPostBlinds() {
	time.Sleep(5 * time.Second)
	e.transitionState(StateDealCards)
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
	e.transitionState(StateProcessGameCommands)
}

func (e *engine) pauseAfterEveryoneFolded() {
	time.Sleep(PAUSE_MEDIUM)
	e.transitionState(StateEveryoneFoldedPayout)
}

func (e *engine) everyoneFoldedPayout() {
	winner := e.state.psuedoDealer

	chips := 0.0
	for _, player := range e.state.players {
		chips += player.chipsInPot
		player.chipsInPot = 0
	}
	chips += e.state.pot

	winner.chips += chips
	e.transitionState(StatePauseAfterEveryoneFoldedPayout)
}

func (e *engine) pauseAfterEveryoneFoldedPayout() {
	time.Sleep(PAUSE_MEDIUM)
	e.transitionState(StateEndHand)
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
	e.state.resetState()
	e.processSitCommand()
	e.transitionState(StateStartHand)
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
