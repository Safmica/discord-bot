package controllers

import (
	"fmt"
	"time"

	models "github.com/Safmica/discord-bot/features/werewolf/models"
	"github.com/bwmarrin/discordgo"
)

var gameStatus = true
var playerReady = 0
var nightNumber = 1
var wwTextID =""
var werewolfEatStatus = false
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
		if voteStatus {
			voteLock.Unlock()
			break
		}
		voteLock.Unlock()

		time.Sleep(1 * time.Second)
	}
}
