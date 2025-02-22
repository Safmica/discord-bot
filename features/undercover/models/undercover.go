package models

import "sync"

type Role string

const (
	Civilian   Role = "civilian"
	Undercover Role = "undercover"
	MrWhite    Role = "mr_white"
)

type WordEntry struct {
	Word []string `json:"word"`
	Used bool     `json:"used"`
}

type WordData struct {
	Words []WordEntry `json:"words"`
}

type Player struct {
	ID       string
	Username string
	Role     Role
}

type GameSession struct {
	ID              string
	Players         map[string]*Player
	Started         bool
	HostID          string
	GameMessageID   string
	Mutex           sync.Mutex
	Undercover      int
	ShowRoles       bool
	Words           []WordEntry
	CivilianWords   string
	UndercoverWords string
}

var ActiveGame *GameSession
