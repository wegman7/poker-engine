package main

// import (
// 	"encoding/json"
// )

type state struct {
	bigBlind float64
	timebankTotal float64
	players map[int]player
	spotlight *player
	spotlightBetweenHands *player
	betweenHands bool
}

func createState(bigBlind float64, timebankTotal float64) *state {
	return &state{
		bigBlind: bigBlind,
		timebankTotal: timebankTotal,
		players: make(map[int]player),
		spotlight: nil,
		spotlightBetweenHands: nil,
		betweenHands: true,
	}
}