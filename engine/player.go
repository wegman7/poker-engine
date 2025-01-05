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
	maxWin		  	float64
	timeBank        float64
	holeCards       []poker.Card
	commandHandlers map[string]commandHandler
	nextInHand      *player
	next            *player
}

type commandHandler func(event *Event, e *engine, s *state) error

func createPlayer(event Event) *player {
	p := player{
		seatId:       event.SeatId,
		user:         event.User,
		sittingOut:   false,
		chips:        event.Chips,
		chipsInPot:   0.0,
		maxWin:       0.0,
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


func (p *player) makeAction(event *Event, e *engine, s *state) error {
	err := p.commandHandlers[event.EngineCommand](event, e, s)
	return err
}

// Add chips to the player's total
func (p *player) addChips(event *Event, e *engine, s *state) error {
	p.chips = p.chips + event.Chips
	return nil
}

// Add chips to the player's total
func (p *player) sitOut(event *Event, e *engine, s *state) error {
	p.sittingOut = true
	return nil
}

// Add chips to the player's total
func (p *player) sitIn(event *Event, e *engine, s *state) error {
	p.sittingOut = false
	return nil
}

// Add chips to the player's total
func (p *player) fold(event *Event, e *engine, s *state) error {
	if err := p.verifySpotlight(s); err != nil {
		return err
	}

	s.removePlayerInHand(p)
	if s.isEveryoneFolded() {
		e.transitionState(StatePauseAfterEveryoneFolded)
	}
	return nil
}

// Add chips to the player's total
func (p *player) check(event *Event, e *engine, s *state) error {
	if err := p.verifySpotlight(s); err != nil {
		return err
	}

	s.rotateSpotlight()
	if s.isStreetComplete() {
		e.transitionState(StateEndStreet)
	}

	return nil
}

func (p *player) call(event *Event, e *engine, s *state) error {
	if err := p.verifySpotlight(s); err != nil {
		return err
	}
	if err := p.verifyLegalCall(s, event.Chips); err != nil {
		return err
	}

	amount := min(s.currentBet - p.chipsInPot, p.chips)
	p.putChipsInPot(s, amount)

	s.rotateSpotlight()
	if s.isStreetComplete() {
		e.transitionState(StateEndStreet)
	}

	return nil
}

// Add chips to the player's total
func (p *player) bet(event *Event, e *engine, s *state) error {
	if err := p.verifySpotlight(s); err != nil {
		return err
	}
	
	betAmount := min(event.Chips - p.chipsInPot, p.chips)
	if err := p.verifyLegalBet(s, betAmount); err != nil {
		return err
	}

	p.putChipsInPot(s, betAmount)

	s.minRaise = betAmount
	s.lastAggressor = p
	s.currentBet = p.chipsInPot
	s.rotateSpotlight()

	return nil
}

func (p *player) putChipsInPot(s *state, amount float64) {
	s.pot += amount
	p.chipsInPot += amount
	p.chips -= amount
}

func (p *player) shouldCreateSidePot(amount float64) bool {
	return amount > p.chipsInPot + p.chips
}

func (p *player) isAllIn() bool {
	return p.chips == 0
}

func (p *player) verifySpotlight(s *state) error {
	if s.spotlight != p {
		return errors.New("it is not your turn")
	}
	return nil
}

func (p *player) verifyLegalCall(s *state, betAmount float64) error {
	if p.chipsInPot == s.currentBet {
		return errors.New("player has already matched the bet")
	}

	return nil
}

func (p *player) verifyLegalBet(s *state, betAmount float64) error {
	// if betAmount == p.chips then the player is all in and the bet is legal
	if betAmount < s.minRaise && betAmount != p.chips {
		return errors.New("bet amount is less than minimum")
	}

	return nil
}

func comparePlayers(prev *player, curr *player) bool {
	if !CompareCardSlices(prev.holeCards, curr.holeCards) {
		return false
	}

	// Compare other fields
	return prev.sittingOut == curr.sittingOut &&
		prev.chips == curr.chips &&
		prev.chipsInPot == curr.chipsInPot &&
		prev.timeBank == curr.timeBank
}
