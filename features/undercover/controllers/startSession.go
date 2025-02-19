package controllers

import (
	"fmt"
	"math/rand"

	models "github.com/Safmica/discord-bot/features/undercover"
	"github.com/bwmarrin/discordgo"
)

func StartGameSession(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if models.ActiveGame == nil || models.ActiveGame.Started {
		return
	}

	playerCount := len(models.ActiveGame.Players)
	if playerCount < 3 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ Minimal 4 pemain diperlukan untuk memulai game.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	models.ActiveGame.Started = true

	players := make([]*models.Player, 0, playerCount)
	for _, p := range models.ActiveGame.Players {
		players = append(players, p)
	}

	rand.Shuffle(len(players), func(i, j int) { players[i], players[j] = players[j], players[i] })

	civilianWord := "Apple"
	undercoverWord := "Orange"

	players[0].Role = models.Undercover
	for i := 1; i < len(players); i++ {
		players[i].Role = models.Civilian
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "🚀 Game telah dimulai! Peranmu akan dikirim secara rahasia melalui DM.",
		},
	})

	for _, p := range players {
		word := ""
		switch p.Role {
		case models.Civilian:
			word = civilianWord
		case models.Undercover:
			word = undercoverWord
		}

		dmChannel, err := s.UserChannelCreate(p.ID)
		if err != nil {
			fmt.Println("Gagal membuat DM channel:", err)
			continue
		}

		_, err = s.ChannelMessageSend(dmChannel.ID, fmt.Sprintf("🔐 **Kata Rahasia Kamu:** %s", word))
		if err != nil {
			fmt.Println("Gagal mengirim DM:", err)
		}
	}

	s.ChannelMessageSend(i.Interaction.ChannelID, "🎮 Game telah dimulai! Diskusikan dan temukan Undercover!")
}
