package engine

import (
	"errors"
	"fmt"

	"github.com/chehsunliu/poker"
)

type state struct {
	smallBlind            float64
	bigBlind              float64
	timebankTotal         float64
	players               map[int]*player
	spotlight             *player
	spotlightBetweenHands *player
	handInAction          bool
	dealer                *player
	deck				  *poker.Deck
	communityCards		  []poker.Card
}

func createState(smallBlind float64, bigBlind float64, timebankTotal float64) *state {
	return &state{
		smallBlind:            smallBlind,
		bigBlind:              bigBlind,
		timebankTotal:         timebankTotal,
		players:               make(map[int]*player),
		spotlight:             nil,
		spotlightBetweenHands: nil,
		handInAction:          true,
		dealer:                nil,
		deck:				   nil,
		communityCards:		   nil,
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

func (s *state) resetPlayersInHand() {
	pointer := s.dealer
	for {
		pointer.nextInHand = nil
		pointer = pointer.next
		if pointer == s.dealer {
			return
		}
	}
}

func (s *state) validateMinimumPlayersInHand() error {
	if s.dealer == nil {
		return errors.New("dealer is nil")
	}

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

	if count < 2 {
		return errors.New("not enough players in hand")
	}
	return nil
}

func (s *state) countPlayersInHand() int {
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
    if err := s.validateMinimumPlayersInHand(); err != nil {
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
	if s.dealer == nil || s.dealer.sittingOut || s.dealer.nextInHand == nil {
		return "No players"
	}

	result := ""
	pointer := s.dealer
	for {
		result += fmt.Sprint(pointer.seatId) + " -> "
		pointer = pointer.nextInHand
		if pointer == s.dealer {
			result += fmt.Sprint(pointer.seatId)
			break
		}
	}
	return result
}