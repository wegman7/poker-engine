package engine

import (
	"errors"
	"fmt"

	"github.com/chehsunliu/poker"
)

type street int

const (
	BetweenHands street = iota
	Preflop
	Flop
	Turn
	River
)

type state struct {
	smallBlind     float64
	bigBlind       float64
	timebankTotal  float64
	players        map[string]*player
	spotlight      *player
	dealer         *player
	psuedoDealer   *player
	lastAggressor  *player
	street         street
	pot            float64
	deck           *poker.Deck
	communityCards []poker.Card
}

func createState(smallBlind float64, bigBlind float64, timebankTotal float64) *state {
	return &state{
		smallBlind:     smallBlind,
		bigBlind:       bigBlind,
		timebankTotal:  timebankTotal,
		players:        make(map[string]*player),
		spotlight:      nil,
		dealer:         nil,
		psuedoDealer:   nil,
		lastAggressor:  nil,
		street:         BetweenHands,
		pot:            0.0,
		deck:           nil,
		communityCards: nil,
	}
}

func (s *state) addPlayer(p *player) {
	s.players[p.user] = p
	if s.dealer == nil {
		s.dealer = p
		p.next = p
		return
	}
	pointer := s.dealer
	for {
		// if we've reached the highest seat id in the middle of the circle (and the player to add has the highest or lowest seat id)
		isReachedHighestSeat := pointer.seatId > pointer.next.seatId && (p.seatId > pointer.seatId || p.seatId < pointer.next.seatId)
		// if the seat id is between two players currently sitting
		isBetweenTwoSeats := pointer.seatId < p.seatId && p.seatId < pointer.next.seatId
		// if the seat id is greater than or lessa than any player currently sitting
		isReachedLastSeat := pointer.next == s.dealer

		if isReachedHighestSeat || isBetweenTwoSeats || isReachedLastSeat {
			p.next = pointer.next
			pointer.next = p
			return
		}
		pointer = pointer.next
	}
}

func (s *state) removePlayer(p *player) {
	delete(s.players, p.user)
	if s.dealer.next == s.dealer {
		s.dealer = nil
		return
	}
	pointer := s.dealer
	for {
		if pointer.next == p {
			pointer.next = p.next
			if s.dealer == p {
				s.dealer = p.next
			}
			break
		}
		pointer = pointer.next
	}
}

// move the psuedo dealer to the next previous in hand
func (s *state) movePsuedoDealer() {
	pointer := s.psuedoDealer
	for {
		if pointer.nextInHand == s.psuedoDealer {
			s.psuedoDealer = pointer
			return
		}
		pointer = pointer.nextInHand
	}
}

func (s *state) removePlayerInHand(p *player) {
	if s.psuedoDealer == p {
		s.movePsuedoDealer()
	}

	pointer := p
	for {
		if pointer.nextInHand == p {
			pointer.nextInHand = p.nextInHand
			p.nextInHand = nil
			return
		}
		pointer = pointer.nextInHand
	}
}

func (s *state) resetDeck() {
	s.deck = nil
	for _, player := range s.players {
		player.holeCards = nil
	}
}

func (s *state) resetPlayers() {
	pointer := s.dealer
	for {
		pointer.nextInHand = nil
		pointer.holeCards = nil

		pointer = pointer.next
		if pointer == s.dealer {
			return
		}
	}
}

func (s *state) resetState() {
	s.resetPlayers()
	s.resetDeck()
	s.spotlight = nil
}

func (s *state) sitoutBustedPlayers() error {
	if s.dealer == nil {
		return errors.New("dealer is nil")
	}

	pointer := s.dealer
	for {
		if pointer.chips == 0 {
			pointer.sittingOut = true
		}
		pointer = pointer.next
		if pointer == s.dealer {
			return nil
		}
	}
}

func (s *state) validateMinimumPlayersSittingIn() error {
	if s.dealer == nil {
		return errors.New("dealer is nil")
	}

	count := s.countPlayersSittingIn()

	if count < 2 {
		return errors.New("not enough players in hand")
	}
	return nil
}

func (s *state) countPlayersSittingIn() int {
	count := 0
	pointer := s.dealer
	for {
		if !pointer.sittingOut {
			count++
		}
		pointer = pointer.next
		if pointer == s.dealer {
			break
		}
	}
	return count
}

func (s *state) countPlayersInHand() int {
	count := 0
	if s.psuedoDealer == nil {
		return count
	}

	pointer := s.psuedoDealer
	for {
		count++
		pointer = pointer.nextInHand
		if pointer == s.psuedoDealer {
			break
		}
	}
	return count
}

func (s *state) isEveryoneFolded() bool {
	return s.countPlayersInHand() == 1
}

func (s *state) isStreetComplete() bool {
	return s.spotlight == s.lastAggressor
}

func (s *state) rotateDealer() error {
	if s.dealer == nil {
		return errors.New("dealer is nil")
	}

	pointer := s.dealer
	for {
		pointer = pointer.next
		if pointer == s.dealer {
			return errors.New("not enough players in hand")
		}
		if !pointer.sittingOut {
			s.dealer = pointer
			s.psuedoDealer = pointer
			return nil
		}
	}
}

func (s *state) orderPlayersInHand() error {
	if s.dealer.sittingOut {
		return errors.New("dealer is sitting out")
	}

	pointer1 := s.dealer
	pointer2 := s.dealer.next
	for {
		for pointer2.sittingOut {
			pointer2 = pointer2.next
		}
		if pointer2 == s.dealer {
			pointer1.nextInHand = pointer2
			break
		}
		pointer1.nextInHand = pointer2
		pointer1 = pointer2
		pointer2 = pointer2.next
	}

	if s.dealer.nextInHand == s.dealer {
		return errors.New("not enough players in hand")
	}
	return nil
}

func (s *state) performDealerRotation() error {
	if err := s.sitoutBustedPlayers(); err != nil {
		return err
	}

	if err := s.validateMinimumPlayersSittingIn(); err != nil {
		return err
	}

	if err := s.rotateDealer(); err != nil {
		return err
	}

	if err := s.orderPlayersInHand(); err != nil {
		return err
	}

	return nil
}

func (s *state) printPlayers() string {
	if s.dealer == nil {
		return "No players"
	}
	result := ""
	pointer := s.dealer
	for {
		result += fmt.Sprint(pointer.seatId) + " -> "
		pointer = pointer.next
		if pointer == s.dealer {
			result += fmt.Sprint(pointer.seatId)
			break
		}
	}
	return result
}

func (s *state) printPlayersInHand() string {
	result := ""
	pointer := s.psuedoDealer
	for {
		result += fmt.Sprint(pointer.seatId) + " -> "
		pointer = pointer.nextInHand
		if pointer == s.psuedoDealer {
			result += fmt.Sprint(pointer.seatId)
			break
		}
	}
	return result
}
