package main

// import (
// 	"encoding/json"
// )

type state struct {
	smallBlind float64
	bigBlind float64
	timebankTotal float64
	players map[int]*player
	spotlight *player
	spotlightBetweenHands *player
	betweenHands bool
	dealer *player
}

func createState(smallBlind float64, bigBlind float64, timebankTotal float64) *state {
	return &state{
		smallBlind: smallBlind,
		bigBlind: bigBlind,
		timebankTotal: timebankTotal,
		players: make(map[int]*player),
		spotlight: nil,
		spotlightBetweenHands: nil,
		betweenHands: true,
		dealer: nil,
	}
}

func (s *state) addPlayer(p *player) {
	s.players[p.seatId] = p
	if s.dealer == nil {
		s.dealer = p
		p.nextPlayerOutOfHand = p
	}
	pointer := s.dealer
	for {
		// if seat id is between to players currently sitting
		if pointer.seatId < p.seatId && pointer.nextPlayerOutOfHand.seatId > p.seatId {
			p.nextPlayerOutOfHand = pointer.nextPlayerOutOfHand
			pointer.nextPlayerOutOfHand = p
			break
		// if the seat id is less than any player currently sitting
		} else if pointer.nextPlayerOutOfHand == s.dealer && p.seatId < pointer.seatId {
			p.nextPlayerOutOfHand = s.dealer
			s.dealer = p
			break
		// if the seat id is greater than any player currently sitting
		} else if pointer.nextPlayerOutOfHand == s.dealer && p.seatId > pointer.seatId {
			pointer.nextPlayerOutOfHand = p
			p.nextPlayerOutOfHand = s.dealer
			s.dealer = p
			break
		}
	}
}

func (s *state) removePlayer(p *player) {
	delete(s.players, p.seatId)
	if s.dealer.nextPlayerOutOfHand == s.dealer {
		s.dealer = nil
		return
	}
	pointer := s.dealer
	for {
		if pointer.nextPlayerOutOfHand == p {
			pointer.nextPlayerOutOfHand = p.nextPlayerOutOfHand
			break
		}
		pointer = pointer.nextPlayerOutOfHand
	}
}

// player (3)
// dealer (0) -> player1 (1) -> player2 (5) -> player3 (6) -> dealer (0)
// dealer (0) -> player1 (1) -> player2 (2) -> dealer (0)
// dealer (4) -> player1 (5) -> player2 (6) -> dealer (4)