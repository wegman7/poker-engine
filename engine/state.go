package engine

import "fmt"

type state struct {
	smallBlind            float64
	bigBlind              float64
	timebankTotal         float64
	players               map[int]*player
	spotlight             *player
	spotlightBetweenHands *player
	betweenHands          bool
	dealer                *player
}

func createState(smallBlind float64, bigBlind float64, timebankTotal float64) *state {
	return &state{
		smallBlind:            smallBlind,
		bigBlind:              bigBlind,
		timebankTotal:         timebankTotal,
		players:               make(map[int]*player),
		spotlight:             nil,
		spotlightBetweenHands: nil,
		betweenHands:          true,
		dealer:                nil,
	}
}

func (s *state) addPlayer(p *player) {
	s.players[p.seatId] = p
	if s.dealer == nil {
		s.dealer = p
		p.next = p
		return
	}
	pointer := s.dealer
	for {
		// if seat id is between two players currently sitting
		if pointer.seatId < p.seatId && pointer.next.seatId > p.seatId {
			p.next = pointer.next
			pointer.next = p
			break
			// if the seat id is greater than or lessa than any player currently sitting
		} else if pointer.next == s.dealer {
			p.next = pointer.next
			pointer.next = p
			break
		}
		pointer = pointer.next
	}
}

func (s *state) removePlayer(p *player) {
	delete(s.players, p.seatId)
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