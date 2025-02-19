package controllers

import (
	models "github.com/Safmica/discord-bot/features/undercover"
	"github.com/bwmarrin/discordgo"
)

func StartGame(s *discordgo.Session, m *discordgo.MessageCreate, i *discordgo.InteractionCreate) {
    if models.ActiveGame != nil {
        sendMessage(s, m, i, "🚀 Game sudah dimulai! Gunakan tombol 'Join Game' untuk bergabung.")
        return
    }

    models.ActiveGame = &models.GameSession{
        ID:      getChannelID(m, i),
        Players: make(map[string]*models.Player),
        Started: false,
    }

    sendMessageWithButton(s, m, i, "🎮 Game Undercover telah dimulai! Klik tombol di bawah untuk bergabung.", "join_game", "Join Game")
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

func sendMessageWithButton(s *discordgo.Session, m *discordgo.MessageCreate, i *discordgo.InteractionCreate, content, buttonID, buttonLabel string) {
    msg := &discordgo.MessageSend{
        Content: content,
        Components: []discordgo.MessageComponent{
            discordgo.ActionsRow{
                Components: []discordgo.MessageComponent{
                    discordgo.Button{
                        Label:    buttonLabel,
                        Style:    discordgo.PrimaryButton,
                        CustomID: buttonID,
                    },
                },
            },
        },
    }

    if m != nil {
        s.ChannelMessageSendComplex(m.ChannelID, msg)
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
