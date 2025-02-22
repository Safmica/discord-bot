package controllers

import (
	"strings"

	models "github.com/Safmica/discord-bot/features/jackheart/models"
	"github.com/bwmarrin/discordgo"
)

func StartGame(s *discordgo.Session, m *discordgo.MessageCreate, i *discordgo.InteractionCreate) {
    if models.ActiveGame != nil {
        sendMessage(s, m, i, "🚀 Game **Jackheart**sudah dimulai! Gunakan tombol 'Join Game' untuk bergabung.")
        return
    }

    hostID := getUserID(m, i)

    models.ActiveGame = &models.GameSession{
        ID:      getChannelID(m, i),
        Players: make(map[string]*models.Player),
        HostID:  hostID, 
        Started: false,
    }

    content :=  "🎮 **Game Jackheart telah dimulai! Klik tombol di bawah untuk bergabung**"

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

func sendMessage(s *discordgo.Session, m *discordgo.MessageCreate, i *discordgo.InteractionCreate, content string){
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
                        CustomID: "join_game",
                    },
                    discordgo.Button{
                        Label:    "Start Game",
                        Style:    discordgo.SuccessButton,
                        CustomID: "start_game",
                    },
                    discordgo.Button{
                        Label:    "Quit Game",
                        Style:    discordgo.DangerButton,
                        CustomID: "quit_game",
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
}

func JackheartHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
    switch i.Type {
    case discordgo.InteractionMessageComponent:
        data := i.MessageComponentData()
        if data.CustomID == "join_game" {
            JoinGame(s, i)
        }

		if data.CustomID == "start_game" {
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

        if data.CustomID == "vote_symbol" {
            ShowVote(s, i, data.CustomID)
        }

        if strings.HasPrefix(data.CustomID, "vote_") {
            HandleVote(s, i, data.CustomID)
        }

        if data.CustomID == "quit_game" {
            QuitGame(s, i)
        }

        if data.CustomID == "play_again" {
            Playagain(s, i)
        }
    }
}
