package engine

import (
	"errors"
	"fmt"
	"math"
	"sort"

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
	collectedPot   float64
	currentBet     float64
	minRaise	   float64
	deck           *poker.Deck
	communityCards []poker.Card
	prevState      *state
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
		collectedPot:   0.0,
		currentBet:     0.0,
		minRaise:		0.0,
		deck:           nil,
		communityCards: nil,
		prevState:      nil,
	}
}

func (s *state) copy() *state {
    copiedPlayers := make(map[string]*player, len(s.players))
    for user, p := range s.players {
        copiedPlayers[user] = p.copy()
    }

    return &state{
        smallBlind:     s.smallBlind,
        bigBlind:       s.bigBlind,
        timebankTotal:  s.timebankTotal,
        players:        copiedPlayers,
        spotlight:      s.spotlight,
        dealer:         s.dealer,
        psuedoDealer:   s.psuedoDealer,
        lastAggressor:  s.lastAggressor,
        street:         s.street,
        pot:            s.pot,
        communityCards: append([]poker.Card{}, s.communityCards...),
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

func (s *state) removePlayersInHand(players []*player) {
	for _, player := range players {
		s.removePlayerInHand(player)
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

// this will keep track of the previous street's pot
func (s *state) collectPot() {
	s.collectedPot = s.pot
}

func (s *state) createSidePots() {
	pointer := s.psuedoDealer
	for {
		if pointer.isAllIn() && pointer.chipsInPot <= s.currentBet {
			pointer.maxWin = s.createSidePot(pointer)
		}
	}
}

func (s *state) createSidePot(hero *player) float64 {
	maxWin := s.collectedPot
	for _, villian := range s.players {
		maxWin += min(hero.chipsInPot, villian.chipsInPot)
	}
	return maxWin
}

func (s *state) isEveryoneFolded() bool {
	return s.countPlayersInHand() == 1
}

func (s *state) isStreetComplete() bool {
	return s.spotlight == s.lastAggressor
}

func (s *state) isStreetRiver() bool {
	return s.street == River
}

func (s *state) isStreetFlop() bool {
	return s.street == Flop
}

func (s *state) goToNextStreet() {
	switch s.street {
	case Preflop:
		s.street = Flop
	case Flop:
		s.street = Turn
	case Turn:
		s.street = River
	}
}

func (s *state) rotateSpotlight() {
	s.spotlight = s.spotlight.nextInHand
	for s.spotlight.isAllIn() && s.spotlight != s.lastAggressor {
		s.spotlight = s.spotlight.nextInHand
	}
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

func (s *state) findBestHand() []*player {
	bestHand := int32(math.MaxInt32)
	winners := make([]*player, 0)

	pointer := s.psuedoDealer
	for pointer != s.psuedoDealer {
		pointer = pointer.nextInHand
		rank := poker.Evaluate(append(pointer.holeCards, s.communityCards...))

		if rank < bestHand {
			winners = make([]*player, 1)
			winners[0] = pointer
		// split pot
		} else if rank == bestHand {
			winners = append(winners, pointer)
		}
		bestHand = min(bestHand, rank)
	}

	return winners
}

func (s *state) payoutWinners(winners []*player) {
	// takes in list of winner(s) and pays them out in order of maxWin asc
	sortWinnersByMaxWin(winners)
	for len(winners) > 0 {
		smallestMaxWin := winners[0]
		if smallestMaxWin.maxWin > 0 {
			s.distributeChips(winners, smallestMaxWin.maxWin)
		}

		winners = winners[1:]
	}
}

// sortWinnersByMaxWin sorts the winners slice by their maxWin in ascending order
func sortWinnersByMaxWin(winners []*player) {
	sort.Slice(winners, func(i, j int) bool {
		return winners[i].maxWin < winners[j].maxWin
	})
}

// distributeChips divides the chips from the smallest maxWin among all winners.
func (s *state) distributeChips(winners []*player, amount float64) {
	for _, winner := range winners {
		winner.chips += amount / float64(len(winners))
		winner.maxWin -= amount
	}
	s.pot -= amount
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

func (s *state) hasStateChanged() bool {
	prev := s.prevState
	curr := s

	if prev == nil {
		s.prevState = curr.copy()
		return true
	}
	s.prevState = curr.copy()

	if prev.dealer != curr.dealer || 
	   prev.psuedoDealer != curr.psuedoDealer || 
	   prev.spotlight != curr.spotlight || 
	   prev.street != curr.street || 
	   prev.pot != curr.pot {
		return true
	}

    if len(prev.players) != len(curr.players) {
        return true
    }

    // compare players
    for user, currPlayer := range curr.players {
        prevPlayer, exists := prev.players[user]
        if !exists || !comparePlayers(prevPlayer, currPlayer) {
            return true
        }
    }

    return false
}