package models

import "sync"

type Role string

const (
	Civilian   Role = "civilian"
	Undercover Role = "undercover"
	MrWhite    Role = "mr_white"
)

type Player struct {
	ID       string
	Username string
	Role     Role
}

type GameSession struct {
	ID            string
	Players       map[string]*Player
	Started       bool
	HostID        string
	GameMessageID string
	Mutex         sync.Mutex
	Undercover    int
	ShowRoles     bool
}

var ActiveGame *GameSession
