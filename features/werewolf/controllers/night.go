package controllers

import (
	"fmt"
	"time"

	models "github.com/Safmica/discord-bot/features/werewolf/models"
	"github.com/bwmarrin/discordgo"
)

var gameStatus = true
var nightNumber = 1
var wwTextID =""
var seerTextID =""
var werewolfEatStatus = false
var seerVoteStatus = false
var lastVoteMessageID = ""
func Night(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// defer func() {
	// 	if r := recover(); r != nil {
	// 		fmt.Println("🚨 Recovered from panic:", r)
	// 	}
	// }()

	if models.ActiveGame == nil || !gameStatus {
		return
	}

	userID := i.Member.User.ID
	channelD := i.ChannelID
	players := models.ActiveGame.Players

	player, exists := models.ActiveGame.Players[userID]
	if !exists || player == nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "⛔ Kamu tidak terdaftar dalam game ini!",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	for gameStatus {
werewolfEatStatus = false
seerVoteStatus = false
	content := fmt.Sprintf(`
**Malam ke %d**
Malam telah tiba, semua pemain kembali kerumah dan beristirahat

		`, nightNumber)

		s.ChannelMessageSend(i.ChannelID, content)
		msg, _ := s.ChannelMessageSendComplex(channelD, &discordgo.MessageSend{
			Content: "_Serigala siap memangsa🐺 (1 menit)_",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label:    "Mangsa Warga",
						Style:    discordgo.PrimaryButton,
						CustomID: "werewolf_eat",
					},
				}}}})

		wwTextID = msg.ID

		for i := 60; i >= 0; i-- {
			voteLock.Lock()
			if !voteWerewolfStatus {
				voteLock.Unlock()
				break
			}
			voteLock.Unlock()

			time.Sleep(1 * time.Second)
		}

		voteWerewolfStatus = false

		if models.ActiveGame.SeerID != "" {
			msg, _ = s.ChannelMessageSendComplex(channelD, &discordgo.MessageSend{
				Content: "_Seer siap meramal 🧙🏻‍♀️ (1 menit)_",
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "Ramal Warga",
							Style:    discordgo.PrimaryButton,
							CustomID: "seer_vote",
						},
					}}}})
			seerTextID = msg.ID
			for i := 60; i >= 0; i-- {
				voteLock.Lock()
				if seerVoteStatus {
					voteLock.Unlock()
					break
				}
				voteLock.Unlock()
		
				time.Sleep(1 * time.Second)
			}
			seerVoteStatus = false
		}

		content = "🌞**Pagi hari telah tiba\n**"
		if len(deathPlayerID) > 0 {
			content += fmt.Sprintf("%d orang telah mati! Warga yang dieliminasi adalah", len(deathPlayerID))
			for _, id := range deathPlayerID {
                content += fmt.Sprintf(" <@%s>,", id)
				delete(models.ActiveGame.Players, id)
            }
		}
		s.ChannelMessageSend(i.ChannelID, content)

		playerList := "📜 **Pemain yang masih hidup:**\n"
		var components []discordgo.MessageComponent
		var buttons []discordgo.MessageComponent

		number := 1
		for _, p := range players {
			playerList += fmt.Sprintf("%d. <@%s>\n", number, p.ID)
			buttons = append(buttons, discordgo.Button{
				Label:    p.Username,
				Style:    discordgo.PrimaryButton,
				CustomID: "werewolf_vote_" + p.ID,
			})
			number++
		}

		if len(buttons) == 5 {
			components = append(components, discordgo.ActionsRow{Components: buttons})
			buttons = []discordgo.MessageComponent{}
		}

		if len(buttons) > 0 {
			components = append(components, discordgo.ActionsRow{Components: buttons})
			buttons = []discordgo.MessageComponent{}
		}

		buttons = append(buttons, discordgo.Button{
			Label:    "Skip",
			Style:    discordgo.SecondaryButton,
			CustomID: "werewolf_vote_skip",
		})

		if len(buttons) == 5 {
			components = append(components, discordgo.ActionsRow{Components: buttons})
			buttons = []discordgo.MessageComponent{}
		}

		if len(buttons) > 0 {
			components = append(components, discordgo.ActionsRow{Components: buttons})
		}

		msg, err := s.ChannelMessageSendComplex(i.Interaction.ChannelID, &discordgo.MessageSend{
			Content:    playerList + "\n🗳 **Silakan pilih pemain yang mencurigakan!**",
			Components: components,
		})

		if err != nil {
			fmt.Println("Gagal mengirim pesan:", err)
		}

		lastVoteMessageID = msg.ID

		for i := 300; i >= 0; i-- {
			voteLock.Lock()
			if villagerVoteStatus {
				voteLock.Unlock()
				break
			}
			voteLock.Unlock()

			time.Sleep(1 * time.Second)
		}
		villagerVoteStatus = false
	}
}
