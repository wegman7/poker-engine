package engine

import (
	"testing"

	"github.com/chehsunliu/poker"
	"github.com/wegman7/game-engine/config"
)

func TestAddRemovePlayer(t *testing.T) {
	s := createState(1, 2, 30)

	p1 := createPlayer(Event{SeatId: 1, User: "user1", Chips: 100})
	p2 := createPlayer(Event{SeatId: 5, User: "user2", Chips: 100})
	p3 := createPlayer(Event{SeatId: 8, User: "user3", Chips: 100})
	p4 := createPlayer(Event{SeatId: 6, User: "user4", Chips: 100})
	p5 := createPlayer(Event{SeatId: 0, User: "user5", Chips: 100})

	s.addPlayer(p1)
	if s.dealer != p1 {
		t.Errorf("Expected dealer to be p1, got %v", s.dealer)
	}
	if s.printPlayers() != "1 -> 1" {
		t.Errorf("Expected 1 -> 1, got %v", s.printPlayers())
	}

	s.addPlayer(p2)
	if s.dealer != p1 {
		t.Errorf("Expected dealer to be p1, got %v", s.dealer)
	}
	if s.printPlayers() != "1 -> 5 -> 1" {
		t.Errorf("Expected 1 -> 5 -> 1, got %v", s.printPlayers())
	}

	s.addPlayer(p3)
	if s.dealer != p1 {
		t.Errorf("Expected dealer to be p1, got %v", s.dealer)
	}
	if s.printPlayers() != "1 -> 5 -> 8 -> 1" {
		t.Errorf("Expected 1 -> 5 -> 8 -> 1, got %v", s.printPlayers())
	}

	s.addPlayer(p4)
	if s.dealer != p1 {
		t.Errorf("Expected dealer to be p1, got %v", s.dealer)
	}
	if s.printPlayers() != "1 -> 5 -> 6 -> 8 -> 1" {
		t.Errorf("Expected 1 -> 5 -> 6 -> 8 -> 1, got %v", s.printPlayers())
	}

	s.addPlayer(p5)
	if s.dealer != p1 {
		t.Errorf("Expected dealer to be p1, got %v", s.dealer)
	}
	if s.printPlayers() != "1 -> 5 -> 6 -> 8 -> 0 -> 1" {
		t.Errorf("Expected 1 -> 5 -> 6 -> 8 -> 0 -> 1, got %v", s.printPlayers())
	}

	s.removePlayer(p1)
	if s.dealer != p2 {
		t.Errorf("Expected dealer to be p2, got %v", s.dealer)
	}
	if s.printPlayers() != "5 -> 6 -> 8 -> 0 -> 5" {
		t.Errorf("Expected 5 -> 6 -> 8 -> 0 -> 5, got %v", s.printPlayers())
	}

	s.removePlayer(p5)
	if s.dealer != p2 {
		t.Errorf("Expected dealer to be p2, got %v", s.dealer)
	}
	if s.printPlayers() != "5 -> 6 -> 8 -> 5" {
		t.Errorf("Expected 5 -> 6 -> 8 -> 5, got %v", s.printPlayers())
	}

	s.removePlayer(p4)
	if s.dealer != p2 {
		t.Errorf("Expected dealer to be p2, got %v", s.dealer)
	}
	if s.printPlayers() != "5 -> 8 -> 5" {
		t.Errorf("Expected 5 -> 8 -> 5, got %v", s.printPlayers())
	}

	s.removePlayer(p2)
	if s.dealer != p3 {
		t.Errorf("Expected dealer to be p2, got %v", s.dealer)
	}
	if s.printPlayers() != "8 -> 8" {
		t.Errorf("Expected 8 -> 8, got %v", s.printPlayers())
	}

	s.removePlayer(p3)
	if s.dealer != nil {
		t.Errorf("Expected dealer to be nil, got %v", s.dealer)
	}
}

func TestSitOutBustedPlayers(t *testing.T) {
	s := createState(1, 2, 30)

	p1 := createPlayer(Event{SeatId: 1, User: "user1", Chips: 100})
	p2 := createPlayer(Event{SeatId: 5, User: "user2", Chips: 0})
	p3 := createPlayer(Event{SeatId: 8, User: "user3", Chips: 100})

	s.addPlayer(p1)
	s.addPlayer(p2)
	s.addPlayer(p3)

	err := s.sitoutBustedPlayers()
	if err != nil {
		t.Errorf("Expected nil, got %s", err.Error())
	}

	if p1.sittingOut || !p2.sittingOut || p3.sittingOut {
		t.Errorf("Expected p1 to be sitting out, got %v, %v, %v", p1.sittingOut, p2.sittingOut, p3.sittingOut)
	}
}

func TestVlidateMinimumPlayersSittingIn(t *testing.T) {
	s := createState(1, 2, 30)

	p1 := createPlayer(Event{SeatId: 1, User: "user1", Chips: 100})
	p2 := createPlayer(Event{SeatId: 5, User: "user2", Chips: 100})
	p3 := createPlayer(Event{SeatId: 8, User: "user3", Chips: 100})

	s.addPlayer(p1)
	s.addPlayer(p2)
	s.addPlayer(p3)

	err := s.validateMinimumPlayersSittingIn()
	if err != nil {
		t.Errorf("Expected nil, got %s", err.Error())
	}

	p1.sittingOut = true
	p2.sittingOut = true
	err = s.validateMinimumPlayersSittingIn()
	if err.Error() != "not enough players in hand" {
		t.Errorf("Expected not enough players in hand, got %s", err.Error())
	}

	s.dealer = nil
	err = s.validateMinimumPlayersSittingIn()
	if err.Error() != "dealer is nil" {
		t.Errorf("dealer is nil, got %s", err.Error())
	}
}

func TestRotateDealer(t *testing.T) {
	s := createState(1, 2, 30)

	p1 := createPlayer(Event{SeatId: 1, User: "user1", Chips: 100})
	p2 := createPlayer(Event{SeatId: 5, User: "user2", Chips: 100})
	p3 := createPlayer(Event{SeatId: 8, User: "user3", Chips: 100})

	err := s.rotateDealer()
	if err == nil || err.Error() != "dealer is nil" {
		t.Errorf("Expected dealer is nil")
	}

	s.addPlayer(p1)
	s.rotateDealer()
	if s.dealer.seatId != 1 {
		t.Errorf("Expected 1, got %v", s.dealer.seatId)
	}

	s.addPlayer(p2)
	s.rotateDealer()
	if s.dealer.seatId != 5 {
		t.Errorf("Expected 5, got %v", s.dealer.seatId)
	}
	if s.printPlayers() != "5 -> 1 -> 5" {
		t.Errorf("Expected 5 -> 1 -> 5, got %v", s.printPlayers())
	}

	p1.sittingOut = true
	p2.sittingOut = true
	err2 := s.rotateDealer()
	if err2 == nil || err2.Error() != "not enough players in hand" {
		t.Errorf("Expected not enough players in hand")
	}

	s.addPlayer(p3)
	s.rotateDealer()
	if s.dealer.seatId != 8 {
		t.Errorf("Expected 8, got %v", s.dealer.seatId)
	}
	if s.printPlayers() != "8 -> 1 -> 5 -> 8" {
		t.Errorf("Expected 8 -> 1 -> 5 -> 8, got %v", s.printPlayers())
	}
}

func TestOrderPlayersInHand(t *testing.T) {
	s := createState(1, 2, 30)

	p1 := createPlayer(Event{SeatId: 1, User: "user1", Chips: 100})
	p2 := createPlayer(Event{SeatId: 5, User: "user2", Chips: 100})
	p3 := createPlayer(Event{SeatId: 8, User: "user3", Chips: 100})
	p4 := createPlayer(Event{SeatId: 6, User: "user4", Chips: 100})
	p5 := createPlayer(Event{SeatId: 0, User: "user5", Chips: 100})
	s.addPlayer(p1)
	s.addPlayer(p2)
	s.addPlayer(p3)
	s.addPlayer(p4)
	s.addPlayer(p5)

	p3.sittingOut = true
	s.orderPlayersInHand()
	s.psuedoDealer = s.dealer
	if s.printPlayersInHand() != "1 -> 5 -> 6 -> 0 -> 1" {
		t.Errorf("Expected 1 -> 5 -> 6 -> 0 -> 1, got %v", s.printPlayersInHand())
	}

	p1.sittingOut = true
	p2.sittingOut = true
	s.rotateDealer()
	s.orderPlayersInHand()

	if s.printPlayersInHand() != "6 -> 0 -> 6" {
		t.Errorf("Expected 6 -> 0 -> 6, got %v", s.printPlayersInHand())
	}
}

func TestResetPlayers(t *testing.T) {
	s := createState(1, 2, 30)

	p1 := createPlayer(Event{SeatId: 1, User: "user1", Chips: 100})
	p2 := createPlayer(Event{SeatId: 5, User: "user2", Chips: 100})
	p3 := createPlayer(Event{SeatId: 8, User: "user3", Chips: 100})
	p4 := createPlayer(Event{SeatId: 6, User: "user4", Chips: 100})
	p5 := createPlayer(Event{SeatId: 0, User: "user5", Chips: 100})
	s.addPlayer(p1)
	s.addPlayer(p2)
	s.addPlayer(p3)
	s.addPlayer(p4)
	s.addPlayer(p5)

	s.performDealerRotation()
	s.resetPlayers()

	pointer := s.dealer
	for {
		if pointer.nextInHand != nil {
			t.Errorf("Expected nil, got %v", pointer.nextInHand)
		}
		if pointer == s.dealer {
			return
		}
		pointer = pointer.next
	}
}

func TestRemovePlayerInHand(t *testing.T) {
	s := createState(1, 2, 30)

	p1 := createPlayer(Event{SeatId: 1, User: "user1", Chips: 100})
	p2 := createPlayer(Event{SeatId: 5, User: "user2", Chips: 100})
	p3 := createPlayer(Event{SeatId: 8, User: "user3", Chips: 100})
	s.addPlayer(p1)
	s.addPlayer(p2)
	s.addPlayer(p3)

	s.performDealerRotation()

	s.removePlayerInHand(p2)
	if s.printPlayersInHand() != "1 -> 8 -> 1" {
		t.Errorf("Expected 1 -> 8 -> 1, got %v", s.printPlayersInHand())
	}

	s.removePlayerInHand(p1)
	if s.printPlayersInHand() != "8 -> 8" {
		t.Errorf("Expected 8 -> 8, got %v", s.printPlayersInHand())
	}
}

func TestFindBestHand(t *testing.T) {
	config.DEBUG = false
	communityCards := []poker.Card{
		poker.NewCard("Ah"),
		poker.NewCard("Kh"),
		poker.NewCard("3h"),
		poker.NewCard("6c"),
		poker.NewCard("Ac"),
	}

	p1 := createPlayer(Event{SeatId: 1, User: "user1", Chips: 100})
	p2 := createPlayer(Event{SeatId: 5, User: "user2", Chips: 100})
	p3 := createPlayer(Event{SeatId: 6, User: "user3", Chips: 100})

	p1.nextInHand = p2
	p1.holeCards = []poker.Card{
		poker.NewCard("As"),
		poker.NewCard("Kd"),
	}
	p2.nextInHand = p3
	p2.holeCards = []poker.Card{
		poker.NewCard("Th"),
		poker.NewCard("9h"),
	}
	p3.nextInHand = p1
	p3.holeCards = []poker.Card{
		poker.NewCard("5s"),
		poker.NewCard("Tc"),
	}
	winners := findBestHand(p1, communityCards)
	if len(winners) != 1 || winners[0] != p1 {
		t.Errorf("Expected p1 to win, got %v", winners)
	}

	p1.holeCards = []poker.Card{
		poker.NewCard("6d"),
		poker.NewCard("3s"),
	}
	p2.holeCards = []poker.Card{
		poker.NewCard("6h"),
		poker.NewCard("3d"),
	}
	winners2 := findBestHand(p1, communityCards)
	if len(winners2) != 2 || winners2[0] != p1 || winners2[1] != p2 {
		t.Errorf("Expected p1 and p2 to split, got %v", winners2)
	}

	p3.holeCards = []poker.Card{
		poker.NewCard("As"),
		poker.NewCard("Kd"),
	}
	winners3 := findBestHand(p1, communityCards)
	if len(winners3) != 1 || winners3[0] != p3 {
		t.Errorf("Expected p3 to win, got %v", winners2)
	}
}

func TestPayoutWinners(t *testing.T) {
	s := createState(1, 2, 30)
	s.communityCards = []poker.Card{
		poker.NewCard("Ah"),
		poker.NewCard("Kh"),
		poker.NewCard("3h"),
		poker.NewCard("6c"),
		poker.NewCard("Ac"),
	}
	s.pot = 1900

	p1 := createPlayer(Event{SeatId: 1, User: "user1", Chips: 0})
	p2 := createPlayer(Event{SeatId: 5, User: "user2", Chips: 0})
	p3 := createPlayer(Event{SeatId: 6, User: "user3", Chips: 0})
	p4 := createPlayer(Event{SeatId: 7, User: "user3", Chips: 0})
	s.addPlayer(p1)
	s.addPlayer(p2)
	s.addPlayer(p3)
	s.addPlayer(p4)

	p1.maxWin = 800
	p2.maxWin = 1100
	p3.maxWin = 1300
	p4.maxWin = 1900

	winners := []*player{p1, p2, p3, p4}

	// we only need cards so we can log winning hand without an error, this won't affect the test because winners are already set
	p1.holeCards = []poker.Card{
		poker.NewCard("As"),
		poker.NewCard("Kd"),
	}
	p1.holeCards = []poker.Card{
		poker.NewCard("As"),
		poker.NewCard("Kd"),
	}
	p1.holeCards = []poker.Card{
		poker.NewCard("As"),
		poker.NewCard("Kd"),
	}
	p1.holeCards = []poker.Card{
		poker.NewCard("As"),
		poker.NewCard("Kd"),
	}
	p1.nextInHand = p2
	p2.nextInHand = p3
	p3.nextInHand = p4
	p4.nextInHand = p1
	s.psuedoDealer = p1
	s.payoutWinners(winners)

	if p1.chips != 200 || p2.chips != 300 || p3.chips != 400 || p4.chips != 1000 {
		t.Errorf("Expected 200, 300, 400, 100, got %v, %v, %v, %v", p1.chips, p2.chips, p3.chips, p4.chips)
	}
}

func TestCreatSidepots(t *testing.T) {
	p1 := createPlayer(Event{SeatId: 1, User: "user1", Chips: 0})
	p2 := createPlayer(Event{SeatId: 5, User: "user2", Chips: 0})
	p3 := createPlayer(Event{SeatId: 6, User: "user3", Chips: 0})
	p4 := createPlayer(Event{SeatId: 7, User: "user4", Chips: 100})

	p1.next, p1.nextInHand = p2, p2
	p1.chipsInPot = 100
	p1.chips = 0

	p2.next, p2.nextInHand = p3, p3
	p2.chipsInPot = 200
	p2.chips = 0

	p3.next, p3.nextInHand = p4, p4
	p3.chipsInPot = 300
	p3.chips = 0

	p4.next, p4.nextInHand = p1, p1
	p4.chipsInPot = 400
	p4.chips = 100

	currentBet := 300.0
	collectedPot := 1000.0
	pot := 2000.0

	createSidePots(p1, currentBet, collectedPot, pot)
	if p1.maxWin != 1400 || p2.maxWin != 1700 || p3.maxWin != 1900 || p4.maxWin != 2000 {
		t.Errorf("Expected 1400, 1700, 1900, 2000, got %v, %v, %v, %v", p1.chips, p2.chips, p3.chips, p4.chips)
	}
}