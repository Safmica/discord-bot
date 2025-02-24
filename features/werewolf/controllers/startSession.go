package controllers

import (
	"fmt"
	"math/rand"
	"sync"

	models "github.com/Safmica/discord-bot/features/werewolf/models"
	"github.com/bwmarrin/discordgo"
)

var playerVotes = make(map[string]string)
var voteCount = make(map[string]int)
var voteMessageID string
var voteStatus bool
var voteLock sync.Mutex

func StartGameSession(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if models.ActiveGame == nil || models.ActiveGame.Started {
		return
	}

	playerCount := len(models.ActiveGame.Players)
	if playerCount < 0 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ Minimal 5 pemain diperlukan untuk memulai game.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	if models.ActiveGame.GameMessageID != "" {
		content := "🎠 **Game Telah dimulai!**"
		s.ChannelMessageEditComplex(&discordgo.MessageEdit{
			ID:         models.ActiveGame.GameMessageID,
			Channel:    models.ActiveGame.ID,
			Content:    &content,
			Components: &[]discordgo.MessageComponent{},
		})
	}

	models.ActiveGame.Started = true

	players := make([]*models.Player, 0, playerCount)
	for _, p := range models.ActiveGame.Players {
		players = append(players, p)
	}

	rand.Shuffle(len(players), func(i, j int) { players[i], players[j] = players[j], players[i] })

	j := 0
	for i := 0; i < models.ActiveGame.Werewolf; i++ {
		players[i].Role = models.Werewolf
		j = models.ActiveGame.Werewolf
		werewolfID = append(werewolfID, players[i].ID)
	}

	for i := j; i < j+models.ActiveGame.Seer; i++ {
		players[i].Role = models.Seer
    }

	for i := j+models.ActiveGame.Seer; i < len(players); i++ {
		players[i].Role = models.Villager
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "🎮 **Game Dimulai!** Tekan tombol di bawah untuk melihat kata rahasiamu!",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "View Your Role",
							Style:    discordgo.PrimaryButton,
							CustomID: "view_role",
						},
					},
				},
			},
		},
	})

	Night(s, i)
}

func ViewRole(s *discordgo.Session, i *discordgo.InteractionCreate) {
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

	content:= fmt.Sprintf("🎁Rolemu adalah **%s**", player.Role)

	if player.Role == "werewolf" {
		if models.ActiveGame.Werewolf > 1 {
			content= fmt.Sprintf("🎁Rolemu adalah **%s** Temanmu adalah:\n", player.Role)
			for _, werewolf := range werewolfID {
				if werewolf != userID {
					friend := models.ActiveGame.Players[werewolf]
					content+=fmt.Sprintf("**%s**, sebagai **%s**\n", friend.Username, friend.Role)
				}
			}
		}
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}