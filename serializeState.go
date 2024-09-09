package main

type SerializePlayer struct {
    SeatId string `json:"seatId"`
}

type SerializeState struct {
    ChannelCommand string    `json:"channelCommand"`
    RoomName string    `json:"roomName"`
    Players map[string]SerializePlayer `json:"data"`
}