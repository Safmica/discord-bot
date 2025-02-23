package controllers

import (
	"fmt"

	models "github.com/Safmica/discord-bot/features/jackheart/models"
	"github.com/bwmarrin/discordgo"
)

func Playagain(s *discordgo.Session, i *discordgo.InteractionCreate) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("🚨 Recovered from panic:", r)
		}
	}()
	voteLock.Lock()
	defer voteLock.Unlock()
	
	if models.ActiveGame != nil {
		sendMessage(s, nil, i, "🚀 Game sudah dimulai! Gunakan tombol 'Join Game' untuk bergabung.")
		return
	}

	hostID := getUserID(nil, i)

	models.ActiveGame = &models.GameSession{
		ID:      getChannelID(nil, i),
		Players: make(map[string]*models.Player),
		HostID:  hostID,
		Started: false,
	}

	gameStatus = true

	sendMessageWithButtons(models.ActiveGame, s, nil, i, "🎮 Game Jackheart telah dimulai! Klik tombol di bawah untuk bergabung.")
}
