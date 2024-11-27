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
var PAUSE_SHORT = 1 * time.Millisecond
var PAUSE_MEDIUM = 2 * time.Millisecond
var PAUSE_LONG = 3 * time.Millisecond

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
	StateShowdown
	StatePauseAfterShowdown
	StateEndHand
	StateDealStreet
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
	case StatePauseAfterEveryoneFolded:
		e.pauseAfterEveryoneFolded()
	case StateEndStreet:
		e.endStreet()
	case StatePauseAfterEndStreet:
		e.pauseAfterEndStreet()
	case StateDealStreet:
		e.dealStreet()
	case StateEveryoneFoldedPayout:
		e.everyoneFoldedPayout()
	case StatePauseAfterEveryoneFoldedPayout:
		e.pauseAfterEveryoneFoldedPayout()
	case StateShowdown:
		e.showdown()
	case StatePauseAfterShowdown:
		e.pauseAfterShowdown()
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
			e.state.players[user].makeAction(&command, e, e.state)
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
		err := e.state.players[user].makeAction(&command, e, e.state)
		if err != nil {
			log.Println("Error processing game command: ", err)
		}
	}
}

func (e *engine) startHand() {
	if err := e.state.performDealerRotation(); err != nil {
		log.Println("Error rotating dealer: ", err)
		return
	}
	e.state.street = Preflop
	e.transitionState(StatePauseAfterStartHand)
}

func (e *engine) pauseAfterStartHand() {
	time.Sleep(PAUSE_MEDIUM)
	e.transitionState(StatePostBlinds)
}

func (e *engine) postBlinds() {
	playerCount := e.state.countPlayersInHand()
	var sb *player
	var bb *player
	if playerCount == 2 {
		sb = e.state.dealer
		bb = e.state.dealer.nextInHand
		e.state.spotlight = sb
	} else {
		sb = e.state.dealer.nextInHand
		bb = e.state.dealer.nextInHand.nextInHand
		e.state.spotlight = bb.nextInHand
	}
	sb.putChipsInPot(e.state, e.state.smallBlind)
	bb.putChipsInPot(e.state, e.state.bigBlind)

	e.state.lastAggressor = bb
	e.state.minRaise = e.state.bigBlind
	e.state.currentBet = e.state.bigBlind

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
	e.state.createSidePots()
	e.state.collectPot()
	e.transitionState(StatePauseAfterEndStreet)
}

func (e *engine) pauseAfterEndStreet() {
	time.Sleep(PAUSE_MEDIUM)
	if e.state.isStreetRiver() {
		e.transitionState(StateShowdown)
	} else {
		e.state.goToNextStreet()
		e.transitionState(StateDealStreet)
	}
}

func (e *engine) resetSpotlight() {
	e.state.spotlight = e.state.psuedoDealer.nextInHand
	e.state.lastAggressor = e.state.psuedoDealer
	e.state.minRaise = e.state.bigBlind
}

func (e *engine) dealStreet() {
	var cards []poker.Card
	if e.state.isStreetFlop() {
		cards = append(cards, poker.NewDeck().Draw(3)...)
	}
	e.state.communityCards = append(e.state.communityCards, cards...)

	e.resetSpotlight()
	e.transitionState(StateProcessGameCommands)
}

func (e *engine) showdown() {
	for e.state.pot > 0 {
		winners := e.state.findBestHand()
		e.state.payoutWinners(winners)
		// remove winners in case we still need to payout a side pot
		e.state.removePlayersInHand(winners)
	}
}

func (e *engine) pauseAfterShowdown() {
	time.Sleep(PAUSE_MEDIUM)
	if e.state.pot > 0 {
		e.transitionState(StateShowdown)
	} else {
		e.transitionState(StateEndHand)
	}
}

func (e *engine) endHand() {
	e.state.resetState()
	e.processSitCommand()
	e.transitionState(StateStartHand)
}

func (e *engine) sendState() {
	if !e.state.hasStateChanged() {
		return
	}

	serializeState := createSerializeState(e.state)
	responseMsg, err := json.Marshal(serializeState)
	if err != nil {
		return
	}

	e.conn.WriteMessage(websocket.TextMessage, responseMsg)
	fmt.Println("Sending state...")
}
