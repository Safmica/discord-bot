package controllers

import (
	"fmt"

	models "github.com/Safmica/discord-bot/features/undercover"
	"github.com/bwmarrin/discordgo"
)

func JoinGame(s *discordgo.Session, i *discordgo.InteractionCreate) {
    if models.ActiveGame == nil {
        s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
            Type: discordgo.InteractionResponseChannelMessageWithSource,
            Data: &discordgo.InteractionResponseData{
                Content: "🚫 Tidak ada game yang sedang berjalan. Gunakan `/startgame` untuk memulai.",
                Flags:   discordgo.MessageFlagsEphemeral,
            },
        })
        return
    }

    models.ActiveGame.Mutex.Lock()
    defer models.ActiveGame.Mutex.Unlock()

    userID := i.Interaction.Member.User.ID
    username := i.Interaction.Member.User.Username

    if _, exists := models.ActiveGame.Players[userID]; exists {
        s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
            Type: discordgo.InteractionResponseChannelMessageWithSource,
            Data: &discordgo.InteractionResponseData{
                Content: fmt.Sprintf("✅ %s, kamu sudah bergabung dalam game!", username),
                Flags:   discordgo.MessageFlagsEphemeral,
            },
        })
        return
    }

    models.ActiveGame.Players[userID] = &models.Player{
        ID:       userID,
        Username: username,
        Role:     "",
    }

	updateGameMessage(s, i.ChannelID, i.Message.ID)
}
