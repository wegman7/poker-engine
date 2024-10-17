package engine

import (
	"testing"
)

func TestDealCards(t *testing.T) {
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

	e := &engine{
		state: s,
	}
	s.performDealerRotation()
	e.dealCards()

	pointer := e.state.dealer.nextInHand
	for {
		if len(pointer.holeCards) != 2 {
			t.Error("Expected 2 hole cards, got", len(pointer.holeCards))
		}
		if pointer == e.state.dealer {
			break
		}
		pointer = pointer.nextInHand
	}
}