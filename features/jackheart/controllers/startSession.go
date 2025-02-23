package controllers

import (
	"math/rand"

	models "github.com/Safmica/discord-bot/features/jackheart/models"
	"github.com/bwmarrin/discordgo"
)

func StartGameSession(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if models.ActiveGame == nil || models.ActiveGame.Started {
		return
	}

	playerCount := len(models.ActiveGame.Players)
	if playerCount < 4 {
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

	content := "🎠 **Game Telah dimulai!**"
	s.ChannelMessageEditComplex(&discordgo.MessageEdit{
		ID:         models.ActiveGame.GameMessageID,
		Channel:    models.ActiveGame.ID,
		Content:    &content,
		Components: &[]discordgo.MessageComponent{},
	})

	players := make([]*models.Player, 0, playerCount)
	for _, p := range models.ActiveGame.Players {
		players = append(players, p)
	}

    models.ActiveGame.MaxPoints = (len(players)*2)+(len(players)+2)

	rand.Shuffle(len(players), func(i, j int) { players[i], players[j] = players[j], players[i] })

	for i := 0; i < len(players); i++ {
        if i == 0 {
            players[i].Role = models.Jackheart
			models.ActiveGame.Jackheart = players[i].ID
        } else {
            players[i].Role = models.Pawn
        }
        players[i].Points = len(players)*2
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "🎮 **Game Dimulai!** Tekan Ready di bawah untuk membuka dashboardmu!",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "Ready",
							Style:    discordgo.PrimaryButton,
							CustomID: "jackheart_dashboard",
						},
					},
				},
			},
		},
	})
}