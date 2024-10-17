package engine

import (
	"fmt"
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

func TestValidateMinimumPlayersInHand(t *testing.T) {
    s := createState(1, 2, 30)

    p1 := createPlayer(Event{SeatId: 1, User: "user1", Chips: 100})
    p2 := createPlayer(Event{SeatId: 5, User: "user1", Chips: 100})
    p3 := createPlayer(Event{SeatId: 8, User: "user1", Chips: 100})

    s.addPlayer(p1)
    s.addPlayer(p2)
    s.addPlayer(p3)

    err := s.validateMinimumPlayersInHand()
    if err != nil {
        t.Errorf("Expected nil, got %s", err.Error())
    }

    p1.sittingOut = true
    p2.sittingOut = true
    err = s.validateMinimumPlayersInHand()
    if err.Error() != "not enough players in hand" {
        t.Errorf("Expected not enough players in hand, got %s", err.Error())
    }

    s.dealer = nil
    err = s.validateMinimumPlayersInHand()
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
    if s.printPlayersInHand() != "1 -> 5 -> 6 -> 0 -> 1" {
        t.Errorf("Expected 1 -> 5 -> 6 -> 0 -> 1, got %v", s.printPlayersInHand())
    }
    
    p1.sittingOut = true
    p2.sittingOut = true
    p4.sittingOut = true
    p5.sittingOut = true
    s.rotateDealer()
    err := s.orderPlayersInHand()
    if err == nil || err.Error() != "dealer is sitting out" {
        t.Errorf("Expected dealer is sitting out")
    }

    p4.sittingOut = false
    p5.sittingOut = false
    s.rotateDealer()
    s.orderPlayersInHand()

    if s.printPlayersInHand() != "6 -> 0 -> 6" {
        t.Errorf("Expected 6 -> 0 -> 6, got %v", s.printPlayersInHand())
    }
}

func TestResetPlayersInHand(t *testing.T) {
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
    
    s.orderPlayersInHand()
    fmt.Println(s.printPlayersInHand())
    s.resetPlayersInHand()
    fmt.Println(s.printPlayersInHand())
}