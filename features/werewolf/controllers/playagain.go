package controllers

import (
	models "github.com/Safmica/discord-bot/features/werewolf/models"
	"github.com/bwmarrin/discordgo"
)

func Playagain(s *discordgo.Session, i *discordgo.InteractionCreate) {
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

	sendMessageWithButtons(models.ActiveGame, s, nil, i, "🎮 Game Undercover telah dimulai! Klik tombol di bawah untuk bergabung.")
}
