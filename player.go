package main

import (
	"github.com/chehsunliu/poker"
)

type player struct{
	seatId int
	user string
	sittingOut bool
	chips float64
	chipsInPot float64
	timeBank float64
	holeCards []poker.Card
}

func createPlayer(user string, chips float64) *player {
	return &player{
		seatId: -1,
		user: user,
		sittingOut: false,
		chips: chips,
		chipsInPot: 0,
		timeBank: 0,
		holeCards: nil,
	}
}