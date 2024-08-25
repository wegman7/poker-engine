package main

// import (
// 	"encoding/json"
// )

type state struct {
	players map[string]player
}

func createState() *state {
	return &state{
		players: make(map[string]player),
	}
}