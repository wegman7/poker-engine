package engine

import (
	"errors"
	"fmt"
	"log"
	"math"
	"sort"

	"github.com/chehsunliu/poker"
	"github.com/wegman7/game-engine/config"
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

func determineSeatId(event Event, players map[string]*player) (int, error) {
	openSeats := make(map[int]bool)
	for i := 0; i < config.MAX_PLAYERS; i++ {
		openSeats[i] = true
	}
	for _, player := range players {
		openSeats[player.seatId] = false
	}

	if event.SeatId != -1 && openSeats[event.SeatId] {
		return event.SeatId, nil
	} else if event.SeatId != -1 && !openSeats[event.SeatId] {
		return -1, errors.New("seat is taken")
	}
	
	return getRandomTrueKey(openSeats)
}

func (s *state) addPlayer(p *player) error {
	if _, exists := s.players[p.user]; exists {
		return errors.New("player already at the table")
	}

	s.players[p.user] = p
	if s.dealer == nil {
		s.dealer = p
		p.next = p
		return nil
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
			return nil
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
	s.communityCards = nil
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
	s.psuedoDealer = nil
	s.lastAggressor = nil
	s.street = BetweenHands
	s.currentBet = 0.0
	s.minRaise = 0.0
	s.pot = 0.0
	s.collectedPot = 0.0
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
	s.currentBet = 0.0

	// reset all players' chips (folded players may still have chips in pot)
	pointer := s.dealer
	for {
		pointer.chipsInPot = 0
		pointer = pointer.next
		if pointer == s.dealer {
			break
		}
	}
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

// decrease the maxWin for all players after the chips have been distributed to the winners (since they're can still be payouts left)
func decreaseMaxWin(psuedoDealer *player, amount float64, winnersSet map[*player]bool) {
	pointer := psuedoDealer
	for {
		if _, exists := winnersSet[pointer]; !exists {
			pointer.maxWin -= amount
		}
		pointer = pointer.nextInHand
		if pointer == psuedoDealer {
			break
		}
	}
}

// distributeChips divides the chips from the smallest maxWin among all winners.
func (s *state) distributeChips(winners []*player, amount float64) {
	winnersSet := make(map[*player]bool)
	for _, winner := range winners {
		winner.chips += amount / float64(len(winners))
		winner.maxWin -= amount
		log.Println(winner.user, " wins ", amount/float64(len(winners)), "with", poker.RankString(poker.Evaluate(append(winner.holeCards, s.communityCards...))))
		winnersSet[winner] = true
	}
	s.pot -= amount
	s.collectedPot -= amount

	decreaseMaxWin(s.psuedoDealer, amount, winnersSet)
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
	   prev.pot != curr.pot ||
	   !CompareCardSlices(prev.communityCards, curr.communityCards) {
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

// creates sidepots (maxWin) for each player in the pot
func createSidePots(psuedoDealer *player, currentBet float64, collectedPot float64, pot float64) {
	pointer := psuedoDealer
	for {
		// we need to check if maxWin is 0 to see if it's already been calculated on a previous street
		if pointer.isAllIn() && pointer.chipsInPot <= currentBet && pointer.maxWin == 0 {
			pointer.maxWin = createSidePot(pointer, collectedPot)
		} else if pointer.maxWin == 0 {
			pointer.maxWin = pot
		}
		pointer = pointer.nextInHand
		if pointer == psuedoDealer {
			break
		}
	}
}

// add the collectedPot (pot from previous street) + chipsInPot from each player (if hero can match it)
func createSidePot(hero *player, collectedPot float64) float64 {
	maxWin := collectedPot
	villian := hero
	for {
		maxWin += min(hero.chipsInPot, villian.chipsInPot)
		villian = villian.next
		if villian == hero {
			break
		}
	}
	return maxWin
}

func findDebugBestHand(seatId int32) int32 {
	if seatId < 5 {
		return 1
	} else {
		return seatId
	}
}

func findBestHand(psuedoDealer *player, communityCards []poker.Card) []*player {
	bestHand := int32(math.MaxInt32)
	winners := make([]*player, 0)

	pointer := psuedoDealer
	for {
		var rank int32
		if config.DEBUG {
			rank = findDebugBestHand(int32(pointer.seatId))
		} else {
			rank = poker.Evaluate(append(pointer.holeCards, communityCards...))
		}

		if rank < bestHand {
			winners = make([]*player, 1)
			winners[0] = pointer
			bestHand = rank
		// split pot
		} else if rank == bestHand {
			winners = append(winners, pointer)
		}

		pointer = pointer.nextInHand
		if pointer == psuedoDealer {
			break
		}
	}

	return winners
}