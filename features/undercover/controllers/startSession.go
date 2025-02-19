package controllers

import (
	"fmt"
	"math/rand"
	"strings"

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

	playerList := "📜 **Daftar Pemain:**\n"
	var buttons []discordgo.MessageComponent

	for _, p := range players {
		playerList += fmt.Sprintf("- <@%s>\n", p.ID)
		buttons = append(buttons, discordgo.Button{
			Label:    p.Username,
			Style:    discordgo.PrimaryButton,
			CustomID: "vote_" + p.ID,
		})
	}

	buttons = append(buttons, discordgo.Button{
		Label:    "Skip",
		Style:    discordgo.DangerButton,
		CustomID: "vote_skip",
	})

	s.ChannelMessageSendComplex(i.Interaction.ChannelID, &discordgo.MessageSend{
		Content: playerList + "\n🗳 **Silakan pilih pemain yang mencurigakan!**",
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{Components: buttons},
		},
	})
}

var playerVotes = make(map[string]string)
var voteCount = make(map[string]int)
var voteMessageID string

func HandleVote(s *discordgo.Session, i *discordgo.InteractionCreate, prefix string) {
    userID := i.Member.User.ID
    voteTarget := strings.TrimPrefix(prefix, "vote_")

    if _, voted := playerVotes[userID]; voted {
        s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
            Type: discordgo.InteractionResponseChannelMessageWithSource,
            Data: &discordgo.InteractionResponseData{
                Content: "❌ Kamu sudah memilih! Tidak bisa memilih lagi.",
                Flags:   discordgo.MessageFlagsEphemeral,
            },
        })
        return
    }

    playerVotes[userID] = voteTarget
    voteCount[voteTarget]++

    voteResults := "📊 **Hasil Voting Sementara:**\n"
    for playerID, count := range voteCount {
        voteResults += fmt.Sprintf("- <@%s>: %d suara\n", playerID, count)
    }

    if voteMessageID != "" {
        _, err := s.ChannelMessageEdit(i.Interaction.ChannelID, voteMessageID, voteResults)
        if err != nil {
            fmt.Println("Gagal mengedit pesan voting:", err)
        }
    } else {
        msg, err := s.ChannelMessageSend(i.Interaction.ChannelID, voteResults)
        if err == nil {
            voteMessageID = msg.ID
        }
    }

    s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
        Type: discordgo.InteractionResponseChannelMessageWithSource,
        Data: &discordgo.InteractionResponseData{
            Content: "✅ Vote kamu telah dicatat!",
            Flags:   discordgo.MessageFlagsEphemeral,
        },
    })
}
