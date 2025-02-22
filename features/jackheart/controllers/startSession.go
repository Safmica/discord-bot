package controllers

import (
	"fmt"
	"math/rand"

	models "github.com/Safmica/discord-bot/features/jackheart/models"
	"github.com/bwmarrin/discordgo"
)

func StartGameSession(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if models.ActiveGame == nil || models.ActiveGame.Started {
		return
	}

	playerCount := len(models.ActiveGame.Players)
	if playerCount < 0 {
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

    models.ActiveGame.MaxPoints = (len(players)*2)+(len(players)+2)

	rand.Shuffle(len(players), func(i, j int) { players[i], players[j] = players[j], players[i] })

	for i := 0; i < len(players); i++ {
        if i == 0 {
            players[i].Role = models.Jackheart
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

var lastVoteMessageID string

func SendVotingMessage(s *discordgo.Session, i *discordgo.InteractionCreate, channelID string) {
	if models.ActiveGame == nil || !models.ActiveGame.Started {
		return
	}

	CloseVoting(s, i, channelID)

	playerList := "📜 **Urutan Bermain:**\n"
	var components []discordgo.MessageComponent
	var buttons []discordgo.MessageComponent
	number := 1

	for _, p := range models.ActiveGame.Players {

		playerList += fmt.Sprintf("%d. <@%s>\n", number, p.ID)
		number++

		buttons = append(buttons, discordgo.Button{
			Label:    p.Username,
			Style:    discordgo.PrimaryButton,
			CustomID: "vote_" + p.ID,
		})

		if len(buttons) == 5 {
			components = append(components, discordgo.ActionsRow{Components: buttons})
			buttons = []discordgo.MessageComponent{}
		}
	}

	if len(buttons) > 0 {
		components = append(components, discordgo.ActionsRow{Components: buttons})
		buttons = []discordgo.MessageComponent{}
	}

	buttons = append(buttons, discordgo.Button{
		Label:    "Skip",
		Style:    discordgo.SecondaryButton,
		CustomID: "vote_skip",
	})

	if len(buttons) == 5 {
		components = append(components, discordgo.ActionsRow{Components: buttons})
		buttons = []discordgo.MessageComponent{}
	}

	buttons = append(buttons, discordgo.Button{
		Label:    "Close Vote",
		Style:    discordgo.DangerButton,
		CustomID: "vote_close",
	})

	if len(buttons) > 0 {
		components = append(components, discordgo.ActionsRow{Components: buttons})
	}

	s.ChannelMessageSendComplex(i.Interaction.ChannelID, &discordgo.MessageSend{
		Content:    playerList + "\n🗳 **Silakan pilih aksi berikut!**",
		Components: components,
	})

}

func CloseVoting(s *discordgo.Session, i *discordgo.InteractionCreate, channelID string) {
	if lastVoteMessageID != "" {
		content := "🗳 **Voting telah selesai!**"
		_, err := s.ChannelMessageEditComplex(&discordgo.MessageEdit{
			ID:         lastVoteMessageID,
			Channel:    i.Interaction.ChannelID,
			Content:    &content,
			Components: &[]discordgo.MessageComponent{},
		})
		if err != nil {
			fmt.Println("Gagal menghapus tombol dari pesan voting:", err)
		}
		lastVoteMessageID = ""
	}
}

// func startVoteCountdown(s *discordgo.Session, channelID string, playerID string) {
// 	msg, err := s.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
// 		Content:   fmt.Sprintf("📜 **Silahkan voting simbol <@%s> Dalam waktu 30 detik** \n" ,playerID),
// 	})

//     button,_ := s.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
// 		Content:   fmt.Sprintf("_Selain <@%s> voting bersifat opsional_", playerID),
// 		Components: 
//         []discordgo.MessageComponent{
//             discordgo.Button{
//                 Label:    "Vote Symbol",
//                 Style:    discordgo.PrimaryButton,
//                 CustomID: "vote_symbol",
//             },
//         },
// 	})

// 	if err != nil {
// 		fmt.Println("Gagal mengirim pesan:", err)
// 	}

// 	for i := 29; i >= 0; i-- {
// 		time.Sleep(1 * time.Second)

// 		_, err := s.ChannelMessageEdit(channelID, msg.ID, fmt.Sprintf("📜 **Silahkan voting simbol <@%s> Dalam waktu %d detik**", playerID, i))
// 		if err != nil {
// 			fmt.Println("Gagal mengedit pesan:", err)
// 			return
// 		}
// 	}

// 	time.Sleep(1 * time.Second)
// 	s.ChannelMessageDelete(channelID, msg.ID)
//     s.ChannelMessageDelete(channelID, button.ID)

// 	s.ChannelMessageSend(channelID, "🛑 **Voting telah berakhir!**")
// }

