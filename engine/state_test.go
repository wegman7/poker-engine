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