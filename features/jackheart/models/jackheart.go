package models

import "sync"

type Role string
type Symbol string

const (
	Pawn      Role = "pawn"
	Jackheart Role = "jackheart"
)
const (
	Heart   Symbol = "heart"
	Spade   Symbol = "spade"
	Diamond Symbol = "diamond"
	Club    Symbol = "club"
)

var Symbols = []Symbol{Heart, Spade, Diamond, Club}

type Player struct {
	ID           string
	Username     string
	Role         Role
	Symbol       Symbol
	SymbolVoted  Symbol
	Points       int
	VotingPlayer string
	DashboardID  string
	VoteID       string
	ActionID     string
}

type GameSession struct {
	ID               string
	Players          map[string]*Player
	Started          bool
	HostID           string
	GameMessageID    string
	Jackheart        string
	VotingNow        string
	VotingID         string
	JackVoteID       string
	AnnounceID       string
	TempVoteVotingID string
	MaxPoints        int
	Mutex            sync.Mutex
	NowPlaying       string
}

var ActiveGame *GameSession
