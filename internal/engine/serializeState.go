package engine

import "github.com/chehsunliu/poker"

type SerializePlayer struct {
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
    Pot float64 `json:"pot"`
    CollectedPot float64 `json:"collectedPot"`
    CurrentBet float64 `json:"currentBet"`
    MinRaise float64 `json:"minRaise"`
    CommunityCards []poker.Card `json:"communityCards"`
	Players map[int]SerializePlayer `json:"players"`
    GameStopped bool `json:"gameStopped"`
}

func createSerializeState(s *state, gameStopped bool) SerializeState {
    serializePlayers := make(map[int]SerializePlayer)
    for _, player := range s.players {
        serializePlayers[player.seatId] = createSerializePlayer(player, s)
    }

    return SerializeState{
        ChannelCommand: "sendState",
        BigBlind: s.bigBlind,
        TimebankTotal: s.timebankTotal,
        Pot: s.pot,
        CollectedPot: s.collectedPot,
        CurrentBet: s.currentBet,
        MinRaise: s.minRaise,
        CommunityCards: s.communityCards,
        Players: serializePlayers,
        GameStopped: gameStopped,
    }
}