package bot

import (
	"github.com/bwmarrin/discordgo"
)

func Help(s *discordgo.Session, m *discordgo.MessageCreate, i *discordgo.InteractionCreate) {
	content := "**PILIH OPSI BANTUAN\nCara memulai permainan gunakan $(undercover/jackheart)**"

    sendMessageWithButtons(s, m, i, content)
}

func sendMessageWithButtons(s *discordgo.Session, m *discordgo.MessageCreate, i *discordgo.InteractionCreate, content string) {
    msg := &discordgo.MessageSend{
        Content: content,
        Components: []discordgo.MessageComponent{
            discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label:    "Undercover",
						Style:    discordgo.PrimaryButton,
						CustomID: "undercover_help",
					},
					discordgo.Button{
						Label:    "Jackheart",
						Style:    discordgo.PrimaryButton,
						CustomID: "jackheart_help",
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