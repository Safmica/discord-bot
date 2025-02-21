package controllers

import (
	"fmt"
	"strconv"
	"strings"

	models "github.com/Safmica/discord-bot/features/undercover"
	"github.com/bwmarrin/discordgo"
)

func ConfigUndercover(s *discordgo.Session, m *discordgo.MessageCreate, args string) {
	if models.ActiveGame == nil {
		s.ChannelMessageSend(m.ChannelID, "⛔ Silahkan mulai game terlebih dahulu dengan $undercover")
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
	case "undercover" :
		undercover, err := strconv.Atoi(config[1])
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "⛔ Input tidak valid! Harap masukkan angka.")
			return
		}
	
		if undercover < 0 {
			undercover = -undercover
		}
	
		totalPlayers := len(models.ActiveGame.Players)
		if totalPlayers == 0 {
			s.ChannelMessageSend(m.ChannelID, "⛔ Tidak ada pemain dalam game!")
			return
		}
	
		maxUndercover := totalPlayers / 4 
	
		if undercover > maxUndercover {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("⛔ Jumlah Undercover terlalu banyak! Maksimal %d dari %d pemain.", maxUndercover, totalPlayers))
			return
		}

		if undercover == 0 {
			s.ChannelMessageSend(m.ChannelID, "⛔ Harus ada minimal 1 undercover")
			return
		}
	
		models.ActiveGame.Undercover = undercover
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("✅ Jumlah Undercover diatur menjadi %d.", undercover))
	case "showroles":
		showroles, err := strconv.ParseBool(config[1])
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "⛔ Harus berisi true or false")
			return
		}
	
		models.ActiveGame.ShowRoles = showroles
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("✅ Showroles diatur menjadi %t.", showroles))
	}
}