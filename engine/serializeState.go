package engine

import "github.com/chehsunliu/poker"

type SerializePlayer struct {
    SeatId int `json:"seatId"`
	User string `json:"user"`
	SittingOut bool `json:"sittingOut"`
	Chips float64 `json:"chips"`
	ChipsInPot float64  `json:"chipsInPot"`
	TimeBank float64 `json:"timeBank"`
	HoleCards []poker.Card `json:"holeCards"`
    Spotlight bool `json:"spotlight"`
    Dealer bool `json:"dealer"`
}

func createSerializePlayer(p *player, s *state) SerializePlayer {
    return SerializePlayer{
        SeatId: p.seatId,
        User: p.user,
        SittingOut: p.sittingOut,
        Chips: p.chips,
        ChipsInPot: p.chipsInPot,
        TimeBank: p.timeBank,
        HoleCards: p.holeCards,
        Spotlight: p == s.spotlight,
        Dealer: p == s.dealer,
    }

}

type SerializeState struct {
    ChannelCommand string `json:"channelCommand"`
	BigBlind float64 `json:"bigBlind"`
	TimebankTotal float64 `json:"timebankTotal"`
	Players map[int]SerializePlayer `json:"players"`
}

func createSerializeState(s *state) SerializeState {
    serializePlayers := make(map[int]SerializePlayer)
    for _, player := range s.players {
        serializePlayers[player.seatId] = createSerializePlayer(player, s)
    }

    return SerializeState{
        ChannelCommand: "sendState",
        BigBlind: s.bigBlind,
        TimebankTotal: s.timebankTotal,
        Players: serializePlayers,
    }
}