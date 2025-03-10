package controllers

import (
	"strings"

	models "github.com/Safmica/discord-bot/features/werewolf/models"
	"github.com/bwmarrin/discordgo"
)

func StartGame(s *discordgo.Session, m *discordgo.MessageCreate, i *discordgo.InteractionCreate) {
	if models.ActiveGame != nil {
		sendMessage(s, m, i, "🚀 Game **Werewolf** sudah dimulai! Gunakan tombol 'Join Game' untuk bergabung.")
		return
	}

	hostID := getUserID(m, i)

	models.ActiveGame = &models.GameSession{
		ID:         getChannelID(m, i),
		Players:    make(map[string]*models.Player),
		HostID:     hostID,
		Started:    false,
		ShowRoles:  false,
		Werewolf: 2,
		Seer: 1,
	}

    content := `
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🐺 **WEREWOLF GAMES** 🐺
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🎮 **Game Werewolf telah dimulai! Klik tombol di bawah untuk bergabung**"
`
	sendMessageWithButtons(models.ActiveGame, s, m, i, content)
}

func getUserID(m *discordgo.MessageCreate, i *discordgo.InteractionCreate) string {
	if m != nil {
		return m.Author.ID
	}
	return i.Member.User.ID
}

func getChannelID(m *discordgo.MessageCreate, i *discordgo.InteractionCreate) string {
	if m != nil {
		return m.ChannelID
	}
	return i.ChannelID
}

func sendMessage(s *discordgo.Session, m *discordgo.MessageCreate, i *discordgo.InteractionCreate, content string) {
	if m != nil {
		s.ChannelMessageSend(m.ChannelID, content)
	} else {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: content,
			},
		})
	}
}

func sendMessageWithButtons(game *models.GameSession, s *discordgo.Session, m *discordgo.MessageCreate, i *discordgo.InteractionCreate, content string) {
	msg := &discordgo.MessageSend{
		Content: content,
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label:    "Join Game",
						Style:    discordgo.PrimaryButton,
						CustomID: "join_werewolf",
					},
					discordgo.Button{
						Label:    "Start Game",
						Style:    discordgo.SuccessButton,
						CustomID: "start_werewolf",
					},
					discordgo.Button{
						Label:    "Quit Game",
						Style:    discordgo.DangerButton,
						CustomID: "quit_werewolf",
					},
				},
			},
		},
	}

	if m != nil {
		message, _ := s.ChannelMessageSendComplex(m.ChannelID, msg)
		game.GameMessageID = message.ID
	} else {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content:    content,
				Components: msg.Components,
			},
		})
	}

	models.ActiveGame.Werewolf = 2
	models.ActiveGame.Seer = 1
	models.ActiveGame.ShowRoles = false
}

func WerewolfHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionMessageComponent:
		data := i.MessageComponentData()
		if data.CustomID == "werewolf_help" {
			Help(s, i)
		}

		if data.CustomID == "join_werewolf" {
			JoinGame(s, i)
		}

		if data.CustomID == "seer_vote" {
			StartSeerVoting(s, i, i.ChannelID)
		}

		if strings.HasPrefix(data.CustomID, "seer_vote_") {
			HandleSeerVote(s, i, data.CustomID)
		}

		if strings.HasPrefix(data.CustomID, "werewolf_vote_") {
			HandleVote(s, i, data.CustomID)
		}

		if data.CustomID == "start_werewolf" {
			if models.ActiveGame == nil || models.ActiveGame.HostID != i.Member.User.ID {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "🚫 Hanya pembuat game yang bisa memulai game!",
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				})
				return
			}
			StartGameSession(s, i)
		}

		if strings.HasPrefix(data.CustomID, "werewolf_eat_") {
			HandleWerewolfVote(s, i, data.CustomID)
		}

		if data.CustomID == "quit_werewolf" {
			QuitGame(s, i)
		}

		if data.CustomID == "werewolf_eat" {
			StartWerewolfVoting(s, i, i.ChannelID)
		}

		if data.CustomID == "play_again_werewolf" {
			Playagain(s, i)
		}

		if data.CustomID == "view_role" {
			ViewRole(s, i)
		}
	}
}
