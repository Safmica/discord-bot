package controllers

import (
	"math/rand"

	models "github.com/Safmica/discord-bot/features/undercover"
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

    players := make([]*models.Player, 0, playerCount)
    for _, p := range models.ActiveGame.Players {
        players = append(players, p)
    }

    rand.Shuffle(len(players), func(i, j int) { players[i], players[j] = players[j], players[i] })

    civilianWord := "Apple"
    undercoverWord := "Orange"

    players[0].Role = models.MrWhite
    players[1].Role = models.Undercover
    for i := 2; i < len(players); i++ {
        players[i].Role = models.Civilian
    }

    for _, p := range players {
        word := ""
        switch p.Role {
        case models.Civilian:
            word = civilianWord
        case models.Undercover:
            word = undercoverWord
        case models.MrWhite:
            word = "???"
        }

        s.ChannelMessageSendComplex(p.ID, &discordgo.MessageSend{
            Embed: &discordgo.MessageEmbed{
                Title: "Kata Rahasia Kamu",
                Description: word,
                Color: 0x00ff00,
            },
        })
    }

    s.ChannelMessageSend(i.Interaction.ChannelID, "🚀 Game telah dimulai! Diskusikan dan temukan Undercover & Mr. White!")
}
