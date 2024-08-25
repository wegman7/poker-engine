package main

type SerializePlayer struct {
    SeatId string `json:"seatId"`
}

type SerializeState struct {
    Type string    `json:"type"`
    Players map[string]SerializePlayer `json:"data"`
}