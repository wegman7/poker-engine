package engine

import (
	"testing"
)

func TestAddRemovePlayer(t *testing.T) {
	s := createState(1, 2, 30)

	p1 := createPlayer(Event{SeatId: 1, User: "user1", Chips: 100})
	p2 := createPlayer(Event{SeatId: 5, User: "user1", Chips: 100})
	p3 := createPlayer(Event{SeatId: 8, User: "user1", Chips: 100})
	p4 := createPlayer(Event{SeatId: 6, User: "user1", Chips: 100})
	p5 := createPlayer(Event{SeatId: 0, User: "user1", Chips: 100})

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
	p2 := createPlayer(Event{SeatId: 5, User: "user1", Chips: 0})
	p3 := createPlayer(Event{SeatId: 8, User: "user1", Chips: 100})

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

func TestvalidateMinimumPlayersSittingIn(t *testing.T) {
	s := createState(1, 2, 30)

	p1 := createPlayer(Event{SeatId: 1, User: "user1", Chips: 100})
	p2 := createPlayer(Event{SeatId: 5, User: "user1", Chips: 100})
	p3 := createPlayer(Event{SeatId: 8, User: "user1", Chips: 100})

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
	p2 := createPlayer(Event{SeatId: 5, User: "user1", Chips: 100})
	p3 := createPlayer(Event{SeatId: 8, User: "user1", Chips: 100})

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
	p2 := createPlayer(Event{SeatId: 5, User: "user1", Chips: 100})
	p3 := createPlayer(Event{SeatId: 8, User: "user1", Chips: 100})
	p4 := createPlayer(Event{SeatId: 6, User: "user1", Chips: 100})
	p5 := createPlayer(Event{SeatId: 0, User: "user1", Chips: 100})
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

func TestresetPlayers(t *testing.T) {
	s := createState(1, 2, 30)

	p1 := createPlayer(Event{SeatId: 1, User: "user1", Chips: 100})
	p2 := createPlayer(Event{SeatId: 5, User: "user1", Chips: 100})
	p3 := createPlayer(Event{SeatId: 8, User: "user1", Chips: 100})
	p4 := createPlayer(Event{SeatId: 6, User: "user1", Chips: 100})
	p5 := createPlayer(Event{SeatId: 0, User: "user1", Chips: 100})
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
	p2 := createPlayer(Event{SeatId: 5, User: "user1", Chips: 100})
	p3 := createPlayer(Event{SeatId: 8, User: "user1", Chips: 100})
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
