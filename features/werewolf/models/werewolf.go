package models

import "sync"

type Role string

const (
	Villager Role = "villager"
	Werewolf Role = "werewolf"
	Seer     Role = "seer"
)

type Player struct {
	ID       string
	Username string
	Role     Role
	VotingID string
}

type GameSession struct {
	ID            string
	Players       map[string]*Player
	Started       bool
	HostID        string
	GameMessageID string
	Mutex         sync.Mutex
	Werewolf      int
	ShowRoles     bool
	Seer          int
	SeerID        string
}

var ActiveGame *GameSession
