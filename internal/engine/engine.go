package engine

import (
	"encoding/json"
	"log"
	"time"

	"github.com/chehsunliu/poker"
	"github.com/gorilla/websocket"
	"github.com/wegman7/game-engine/config"
)

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
	StatePauseAfterEndHand
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
			delete(runningEngines, e.roomName)
			log.Println("Stopping engine for room", e.roomName)
			return
		default:
			time.Sleep(config.AppConfig.ENGINE_LOOP_PAUSE)
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
	case StatePauseAfterEndHand:
		e.pauseAfterEndHand()
	}
}

func (e *engine) transitionState(newEngineState engineState) {
	log.Println("Transitioning state from", e.engineState, "to", newEngineState)
	log.Println("Street:", e.state.street, "Community cards:", e.state.communityCards)
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
		log.Println("processing sit command: ", command)
		user := command.User

		if command.EngineCommand == "join" {
			seatId, err := determineSeatId(command, e.state.players)
			if err != nil {
				log.Println("Error determining seat id: ", err)
				continue
			}
			command.SeatId = seatId
			p := createPlayer(command)
			err2 := e.state.addPlayer(p)
			if err2 != nil {
				log.Println("Error adding player: ", err2)
				continue
			}
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
		log.Println("processing game command: ", command)
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
		e.transitionState(StateProcessSitCommands)
		return
	}
	e.state.street = Preflop
	e.transitionState(StatePauseAfterStartHand)
}

func (e *engine) pauseAfterStartHand() {
	time.Sleep(config.AppConfig.PAUSE_MEDIUM)
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
		e.state.lastAggressor = sb
	} else {
		sb = e.state.dealer.nextInHand
		bb = e.state.dealer.nextInHand.nextInHand
		e.state.spotlight = bb.nextInHand
		e.state.lastAggressor = bb.nextInHand
	}
	sb.putChipsInPot(e.state, e.state.smallBlind)
	bb.putChipsInPot(e.state, e.state.bigBlind)

	e.state.minRaise = e.state.bigBlind
	e.state.currentBet = e.state.bigBlind

	e.transitionState(StatePauseAfterPostBlinds)
}

func (e *engine) pauseAfterPostBlinds() {
	time.Sleep(config.AppConfig.PAUSE_MEDIUM)
	e.transitionState(StateDealCards)
}

func (e *engine) dealCards() {
	e.state.deck = poker.NewDeck()
	pointer := e.state.dealer.nextInHand
	for {
		pointer.holeCards = e.state.deck.Draw(2)
		if pointer == e.state.dealer {
			break
		}
		pointer = pointer.nextInHand
	}
	e.transitionState(StateProcessGameCommands)
}

func (e *engine) pauseAfterEveryoneFolded() {
	time.Sleep(config.AppConfig.PAUSE_MEDIUM)
	e.transitionState(StateEveryoneFoldedPayout)
}

func (e *engine) everyoneFoldedPayout() {
	winner := e.state.psuedoDealer
	e.state.collectPot()
	winner.chips += e.state.pot
	e.state.pot -= e.state.pot
	e.state.collectedPot -= e.state.pot
	e.transitionState(StatePauseAfterEveryoneFoldedPayout)
}

func (e *engine) pauseAfterEveryoneFoldedPayout() {
	time.Sleep(config.AppConfig.PAUSE_MEDIUM)
	e.transitionState(StateEndHand)
}

func (e *engine) endStreet() {
	createSidePots(e.state.psuedoDealer, e.state.currentBet, e.state.collectedPot, e.state.pot)
	e.state.collectPot()
	e.transitionState(StatePauseAfterEndStreet)
}

func (e *engine) pauseAfterEndStreet() {
	time.Sleep(config.AppConfig.PAUSE_MEDIUM)
	if e.state.isStreetRiver() {
		e.transitionState(StateShowdown)
	} else {
		e.state.goToNextStreet()
		e.transitionState(StateDealStreet)
	}
}

// returns true if all players are all in
func (e *engine) resetSpotlight() bool {
	e.state.spotlight = e.state.psuedoDealer.nextInHand
	for e.state.spotlight.isAllIn() {
		e.state.spotlight = e.state.spotlight.nextInHand
		if e.state.spotlight == e.state.psuedoDealer.nextInHand {
			return true
		}
	}

	e.state.lastAggressor = e.state.psuedoDealer.nextInHand
	e.state.minRaise = e.state.bigBlind
	return false
}

func (e *engine) dealStreet() {
	var cards []poker.Card
	if e.state.isStreetFlop() {
		cards = append(cards, e.state.deck.Draw(3)...)
	} else {
		cards = append(cards, e.state.deck.Draw(1)...)
	}
	e.state.communityCards = append(e.state.communityCards, cards...)

	// if all players are all in, skip to end street
	isAllPlayersAllIn := e.resetSpotlight()
	if isAllPlayersAllIn {
		e.transitionState(StateEndStreet)
	} else {
		e.transitionState(StateProcessGameCommands)
	}
}

func (e *engine) showdown() {
	winners := findBestHand(e.state.psuedoDealer, e.state.communityCards)
	e.state.payoutWinners(winners)

	// remove winners in case we still need to payout a side pot
	e.state.removePlayersInHand(winners)
	e.transitionState(StatePauseAfterShowdown)
}

func (e *engine) pauseAfterShowdown() {
	time.Sleep(config.AppConfig.PAUSE_MEDIUM)

	// continue to pay sidepots until the pot is empty
	if e.state.pot > 0 {
		e.transitionState(StateShowdown)
	} else {
		e.transitionState(StateEndHand)
	}
}

func (e *engine) endHand() {
	e.state.resetState()
	e.processSitCommand()
	e.transitionState(StatePauseAfterEndHand)
}

func (e *engine) pauseAfterEndHand() {
	time.Sleep(config.AppConfig.PAUSE_LONG)
	e.transitionState(StateStartHand)
}

func (e *engine) sendState() {
	if !e.state.hasStateChanged() {
		return
	}

	serializeState := createSerializeState(e.state, e.engineState == StateProcessSitCommands)
	responseMsg, err := json.Marshal(serializeState)
	if err != nil {
		return
	}

	e.conn.WriteMessage(websocket.TextMessage, responseMsg)
	log.Println("Sending state...")
}
