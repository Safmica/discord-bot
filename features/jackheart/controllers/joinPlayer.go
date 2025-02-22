package controllers

import (
	"fmt"

	models "github.com/Safmica/discord-bot/features/jackheart/models"
	"github.com/bwmarrin/discordgo"
)

func updateGameMessage(s *discordgo.Session, channelID, messageID string) {
	if models.ActiveGame == nil {
		return
	}

	playerList :=  `
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🎭 **JACK HEART DASHBOARD** 🎭
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🎮 **List Pemain**"
`
	for _, player := range models.ActiveGame.Players {
		playerList += fmt.Sprintf("🔹 <@%s>\n", player.ID)
	}

	_, err := s.ChannelMessageEdit(channelID, messageID, playerList)
	if err != nil {
		fmt.Println("Error updating game message:", err)
	}
}
