package controllers

import (
	"fmt"

	models "github.com/Safmica/discord-bot/features/undercover"
	"github.com/bwmarrin/discordgo"
)

func QuitGame(s *discordgo.Session, i *discordgo.InteractionCreate) {
	userID := i.Interaction.Member.User.ID

	models.ActiveGame.Mutex.Lock()
    defer models.ActiveGame.Mutex.Unlock()

	valid := models.ActiveGame.Players[userID]
	if valid != nil {
		delete(models.ActiveGame.Players, userID)
		updateGameMessage(s, i.ChannelID, i.Message.ID)

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("✅ <@%s>, kamu berhasil keluar dari game!", userID),
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
	} else { 
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("⛔ <@%s>, kamu belum pernah bergabung di dalam game!", userID),
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
	}
}