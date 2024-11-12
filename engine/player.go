package engine

import (
	"errors"

	"github.com/chehsunliu/poker"
)

type player struct {
	seatId          int
	user            string
	sittingOut      bool
	chips           float64
	chipsInPot      float64
	timeBank        float64
	holeCards       []poker.Card
	commandHandlers map[string]commandHandler
	nextInHand      *player
	next            *player
}

type commandHandler func(event Event, e *engine, s *state) error

func createPlayer(event Event) *player {
	p := player{
		seatId:       event.SeatId,
		user:         event.User,
		sittingOut:   false,
		chips:        event.Chips,
		chipsInPot:   0,
		timeBank:     0,
		holeCards:    nil,
		nextInHand:   nil,
		next:         nil,
	}

	p.commandHandlers = make(map[string]commandHandler)
	p.commandHandlers["addChips"] = p.addChips
	p.commandHandlers["sitOut"] = p.sitOut
	p.commandHandlers["sitIn"] = p.sitIn
	p.commandHandlers["fold"] = p.fold
	p.commandHandlers["check"] = p.check
	p.commandHandlers["call"] = p.call
	p.commandHandlers["bet"] = p.bet

	return &p
}

func (p *player) copy() *player {
    return &player{
        seatId:      p.seatId,
        user:        p.user,
        sittingOut:  p.sittingOut,
        chips:       p.chips,
        chipsInPot:  p.chipsInPot,
        timeBank:    p.timeBank,
        holeCards:   append([]poker.Card{}, p.holeCards...),
    }
}


func (p *player) makeAction(event Event, e *engine, s *state) error {
	p.commandHandlers[event.EngineCommand](event, e, s)
	return nil
}

// Add chips to the player's total
func (p *player) addChips(event Event, e *engine, s *state) error {
	p.chips = p.chips + event.Chips
	return nil
}

// Add chips to the player's total
func (p *player) sitOut(event Event, e *engine, s *state) error {
	p.sittingOut = true
	return nil
}

// Add chips to the player's total
func (p *player) sitIn(event Event, e *engine, s *state) error {
	p.sittingOut = false
	return nil
}

// Add chips to the player's total
func (p *player) fold(event Event, e *engine, s *state) error {
	if err := p.verifyLegalMove(e, s); err != nil {
		return err
	}

	s.removePlayerInHand(p)
	if s.isEveryoneFolded() {
		e.transitionState(StatePauseAfterEveryoneFolded)
	}
	return nil
}

// Add chips to the player's total
func (p *player) check(event Event, e *engine, s *state) error {
	if err := p.verifyLegalMove(e, s); err != nil {
		return err
	}

	if s.isStreetComplete() {
		e.transitionState(StateEndStreet)
	} else {
		s.spotlight = s.spotlight.nextInHand
	}

	return nil
}

// Add chips to the player's total
func (p *player) call(event Event, e *engine, s *state) error {
	if err := p.verifyLegalMove(e, s); err != nil {
		return err
	}

	return nil
}

// Add chips to the player's total
func (p *player) bet(event Event, e *engine, s *state) error {
	if err := p.verifyLegalMove(e, s); err != nil {
		return err
	}

	return nil
}

func (p *player) postBlind(amount float64) error {
	p.chipsInPot = p.chipsInPot + amount
	p.chips = p.chips - amount
	return nil
}

func (p *player) verifyLegalMove(e *engine, s *state) error {
	if err := p.verifySpotlight(e, s); err != nil {
		return err
	}
	return nil
}

func (p *player) verifySpotlight(e *engine, s *state) error {
	if s.spotlight != p {
		return errors.New("it is not your turn")
	}
	return nil
}