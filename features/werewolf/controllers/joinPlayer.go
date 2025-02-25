package controllers

import (
	"fmt"

	models "github.com/Safmica/discord-bot/features/werewolf/models"
	"github.com/bwmarrin/discordgo"
)

func updateGameMessage(s *discordgo.Session, channelID, messageID string) {
	if models.ActiveGame == nil {
		return
	}

	playerList :=  fmt.Sprintf(`
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🐺 **WEREWOLF DASHBOARD** 🐺
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
_Jumlah Werewolf = %d_ (not realtime)
_Jumlah Seer = %d_ (not realtime)
_Showroles = %t_ (not realtime)

🎮 **List Pemain**"
`, models.ActiveGame.Werewolf,models.ActiveGame.Seer, models.ActiveGame.ShowRoles)
	for _, player := range models.ActiveGame.Players {
		playerList += fmt.Sprintf("🔹 <@%s>\n", player.ID)
	}

	_, err := s.ChannelMessageEdit(channelID, messageID, playerList)
	if err != nil {
		fmt.Println("Error updating game message:", err)
	}
}
