package controllers

import (
	"fmt"

	models "github.com/Safmica/discord-bot/features/undercover"
	"github.com/bwmarrin/discordgo"
)

func updateGameMessage(s *discordgo.Session, channelID, messageID string) {
    if models.ActiveGame == nil {
        return
    }

    playerList := fmt.Sprintf("🎮 **Pemain yang sudah bergabung:**.\n _Jumlah Undercover = %d_ (not realtime)\n _Showroles = %t_ (not realtime)\n", models.ActiveGame.Undercover, models.ActiveGame.ShowRoles)
    for _, player := range models.ActiveGame.Players {
        playerList += fmt.Sprintf("🔹 <@%s>\n", player.ID)
    }

    _, err := s.ChannelMessageEdit(channelID, messageID, playerList)
    if err != nil {
        fmt.Println("Error updating game message:", err)
    }
}
