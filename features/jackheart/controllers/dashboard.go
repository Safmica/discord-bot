package controllers

import (
	"fmt"
	"math/rand"

	models "github.com/Safmica/discord-bot/features/jackheart/models"
	"github.com/bwmarrin/discordgo"
)
var playerReady = 0
func Dashboard(s *discordgo.Session, i *discordgo.InteractionCreate) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("🚨 Recovered from panic:", r)
		}
	}()
	
	if models.ActiveGame == nil || !gameStatus {
		return
	}

	userID := i.Member.User.ID

	player, exists := models.ActiveGame.Players[userID]
	if !exists || player == nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "⛔ Kamu tidak terdaftar dalam game ini!",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	if player.Symbol == "" {
		if len(models.Symbols) == 0 {
			fmt.Println("❌ ERROR: models.Symbols kosong!")
			return
		}
		rand.Shuffle(len(models.Symbols), func(i, j int) { models.Symbols[i], models.Symbols[j] = models.Symbols[j], models.Symbols[i] })
		player.Symbol = models.Symbols[0]
	}

	var role string
	if player.Role == "jackheart" {
		role = "Jackheart 🎭"
	} else {
		role = "Pawn ♟️"
	}

	if player.DashboardID != "" {
		if err := s.FollowupMessageDelete(i.Interaction, player.DashboardID); err != nil {
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

	// ✅ Simpan ID dashboard baru
	player.DashboardID = msg.ID
	playerReady++

	// ✅ Pastikan fase dan jumlah pemain sesuai sebelum lanjut
	if playerReady == len(models.ActiveGame.Players) && Phase == "finish" {
		fmt.Println(player.Symbol)
		playerReady = 0
		startTurnBasedVoting(s, i.ChannelID)
	}
}
