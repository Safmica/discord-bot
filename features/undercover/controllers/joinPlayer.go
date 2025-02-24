package controllers

import (
	"fmt"

	models "github.com/Safmica/discord-bot/features/undercover/models"
	"github.com/bwmarrin/discordgo"
)

func updateGameMessage(s *discordgo.Session, channelID, messageID string) {
	if models.ActiveGame == nil {
		return
	}

	playerList :=  fmt.Sprintf(`
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🎭 **UNDERCOVER DASHBOARD** 🎭
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
_Jumlah Undercover = %d_ (not realtime)
_Jumlah Mr.White = %d_ (not realtime)
_Showroles = %t_ (not realtime)

🎮 **List Pemain**"
`, models.ActiveGame.Undercover,models.ActiveGame.MrWhite, models.ActiveGame.ShowRoles)
	for _, player := range models.ActiveGame.Players {
		playerList += fmt.Sprintf("🔹 <@%s>\n", player.ID)
	}

	_, err := s.ChannelMessageEdit(channelID, messageID, playerList)
	if err != nil {
		fmt.Println("Error updating game message:", err)
	}
}
