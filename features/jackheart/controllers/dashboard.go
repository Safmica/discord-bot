package controllers

import (
	"fmt"
	"math/rand"

	models "github.com/Safmica/discord-bot/features/jackheart/models"
	"github.com/bwmarrin/discordgo"
)
var playerReady = 0
func Dashboard(s *discordgo.Session, i *discordgo.InteractionCreate) {
	userID := i.Member.User.ID
	player, exists := models.ActiveGame.Players[userID]
	if !exists {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "⛔ Kamu tidak terdaftar dalam game ini!",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	rand.Shuffle(len(models.Symbols), func(i, j int) { models.Symbols[i], models.Symbols[j] = models.Symbols[j], models.Symbols[i] })
	player.Symbol = models.Symbols[0]

	var role string
	if player.Role == "jackheart" {
		role = "Jackheart 🎭"
	} else {
		role = "Pawn ♟️"
	}

	if player.DashboardID != "" {
		err := s.FollowupMessageDelete(i.Interaction, player.DashboardID)
		if err != nil {
			fmt.Println("Gagal menghapus pesan lama:", err)
		}
	}

	content := fmt.Sprintf("⚙️ **Role Kamu : %s** \n 🎖️**Point Kamu : %d** _Point Maksimal : %d_", role, player.Points, models.ActiveGame.MaxPoints)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral, 
		},
	})

	msg, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: content,
		Flags:   discordgo.MessageFlagsEphemeral,
	})
	if err != nil {
		fmt.Println("Gagal mengirim pesan follow-up:", err)
		return
	}
	
	player.DashboardID = msg.ID
	playerReady++
	
	if playerReady == len(models.ActiveGame.Players) && Phase == "finish"{
		playerReady = 0
		startTurnBasedVoting(s, i.ChannelID)
	}
}

