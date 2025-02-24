package controllers

import (
	"fmt"
	"strconv"
	"strings"

	models "github.com/Safmica/discord-bot/features/werewolf/models"
	"github.com/bwmarrin/discordgo"
)

func Configwerewolf(s *discordgo.Session, m *discordgo.MessageCreate, args string) {
	if models.ActiveGame == nil {
		s.ChannelMessageSend(m.ChannelID, "⛔ Silahkan mulai game terlebih dahulu dengan $werewolf")
		return
	}

	userID := m.Author.ID
	

	models.ActiveGame.Mutex.Lock()
    defer models.ActiveGame.Mutex.Unlock()

	if models.ActiveGame.HostID != userID {
		s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("⛔ <@%s>, kamu bukan host di game!", m.Author.ID), &discordgo.MessageReference{
			MessageID: m.ID,
		})
		return
	}

	config := strings.Fields(args)
	switch config[0] {
	case "werewolf" :
		werewolf, err := strconv.Atoi(config[1])
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "⛔ Input tidak valid! Harap masukkan angka.")
			return
		}
	
		if werewolf < 0 {
			werewolf = -werewolf
		}
	
		totalPlayers := len(models.ActiveGame.Players)
		if totalPlayers == 0 {
			s.ChannelMessageSend(m.ChannelID, "⛔ Tidak ada pemain dalam game!")
			return
		}
	
		models.ActiveGame.Werewolf = werewolf
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("✅ Jumlah werewolf diatur menjadi %d.", werewolf))
	case "showroles":
		showroles, err := strconv.ParseBool(config[1])
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "⛔ Harus berisi true or false")
			return
		}
	
		models.ActiveGame.ShowRoles = showroles
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("✅ Showroles diatur menjadi %t.", showroles))
	case "seer" :
		seer, err := strconv.Atoi(config[1])
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "⛔ Input tidak valid! Harap masukkan angka.")
			return
		}
	
		if seer < 0 {
			seer = -seer
		}
	
		totalPlayers := len(models.ActiveGame.Players)
		if totalPlayers == 0 {
			s.ChannelMessageSend(m.ChannelID, "⛔ Tidak ada pemain dalam game!")
			return
		}
	
		models.ActiveGame.Seer = seer
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("✅ Jumlah Seer diatur menjadi %d.", seer))
	}
}