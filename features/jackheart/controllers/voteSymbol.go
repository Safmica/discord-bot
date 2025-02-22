package controllers

import (
	"fmt"
	"strings"
	"sync"
	"time"

	models "github.com/Safmica/discord-bot/features/jackheart/models"
	"github.com/bwmarrin/discordgo"
)

var voteMessageID string
var playerVotes = make(map[string]string)
var symbolVotes = make(map[string]string)
var voteStatus bool
var voteLock sync.Mutex

func startTurnBasedVoting(s *discordgo.Session, channelID string) {
	activeGame := models.ActiveGame // Ambil game aktif
	players := activeGame.Players   // Ambil daftar pemain

	for _, player := range players {
		models.ActiveGame.NowPlaying = player.ID
		voteStatus = true // Set awal voting aktif
	
		content := fmt.Sprintf("📜 **Silahkan voting simbol <@%s> Dalam waktu 30 detik**\n _Selain <@%s>, voting bersifat opsional_", player.ID, player.ID)
	
		// Kirim pesan pertama kali
		if voteMessageID == "" {
			msg, _ := s.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
				Content: content,
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label:    "Vote Symbol",
						Style:    discordgo.PrimaryButton,
						CustomID: "vote_symbol",
					},
				},
			})
			voteMessageID = msg.ID
		}
	
		// Countdown voting 30 detik atau sampai player memilih
		for i := 29; i >= 0; i-- {
			if !voteStatus { // Jika pemain sudah vote, hentikan sesi lebih cepat
				break
			}
	
			time.Sleep(1 * time.Second)
	
			voteResults := "📜 **Hasil Voting Sementara:**\n"
			for player, symbol := range playerVotes {
				voteResults += fmt.Sprintf("- <@%s> memilih **%s**\n", player, symbol)
			}
	
		// Gabungkan dengan pesan utama
			content = fmt.Sprintf("📜 **Silahkan voting simbol <@%s> Dalam waktu 30 detik**\n _Selain <@%s>, voting bersifat opsional_\n\n%s", 
			player.ID, player.ID, voteResults)
	
			s.ChannelMessageEditComplex(&discordgo.MessageEdit{
				ID:      voteMessageID,
				Channel: channelID,
				Content: &content,
				Components: &[]discordgo.MessageComponent{
					discordgo.ActionsRow{Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "Vote Symbol",
							Style:    discordgo.PrimaryButton,
							CustomID: "vote_symbol",
						},
					}},
				},
			})
		}
	
		// Hapus pesan voting setelah sesi selesai
		s.ChannelMessageSend(channelID, fmt.Sprintf("🛑 **Voting untuk <@%s> telah berakhir!**", player.ID))
	}
	
	// Setelah semua voting selesai
	s.ChannelMessageSend(channelID, "✅ **Semua pemain telah menyelesaikan voting!**")
}	

func ShowVote(s *discordgo.Session, i *discordgo.InteractionCreate, prefix string) {
	voteLock.Lock()
	defer voteLock.Unlock()

	voterID := i.Member.User.ID
	models.ActiveGame.Players[voterID].Points--
	userID := models.ActiveGame.NowPlaying

	players := models.ActiveGame.Players[userID]
	symbol := ""
	content := ""
	switch players.Symbol {
	case models.Heart:
		symbol = "**Heart ❤️**\n"
	case models.Diamond:
		symbol = "**Diamond ♦️**\n"
	case models.Club:
		symbol = "**Club ♣️**\n"
	case models.Spade:
		symbol = "**Spade ♠️**\n"
	}
	if players.ID == userID {
		content = "📜 **Pilihlah symbol anda**\n"
	} else {
		content = fmt.Sprintf("**Symbol milik <@%s> adalah %s📜 **Pilihlah symbol untuknya**\n_pointmu berkurang -1_", players.ID, symbol)
	}

	var buttons []discordgo.MessageComponent

	buttons = append(buttons, discordgo.Button{
		Label:    "Heart ❤️",
		Style:    discordgo.SecondaryButton,
		CustomID: "vote_heart",
	})

	buttons = append(buttons, discordgo.Button{
		Label:    "Spade ♠️",
		Style:    discordgo.DangerButton,
		CustomID: "vote_spade",
	})

	buttons = append(buttons, discordgo.Button{
		Label:    "Diamond ♦️",
		Style:    discordgo.DangerButton,
		CustomID: "vote_diamond",
	})

	buttons = append(buttons, discordgo.Button{
		Label:    "Club ♣️",
		Style:    discordgo.DangerButton,
		CustomID: "vote_club",
	})

	msg, err := s.ChannelMessageSendComplex(i.Interaction.ChannelID, &discordgo.MessageSend{
		Content:    content,
		Components: buttons,
	})

	if err != nil {
		fmt.Println("Gagal mengirim pesan:", err)
	}

	lastVoteMessageID = msg.ID
}

func HandleVote(s *discordgo.Session, i *discordgo.InteractionCreate, prefix string) {
	if models.ActiveGame == nil || !models.ActiveGame.Started {
		return
	}

	voteStatus = true
	voteLock.Lock()
	defer voteLock.Unlock()
	voterID := i.Member.User.ID
	userID := models.ActiveGame.NowPlaying

	voteTarget := strings.TrimPrefix(prefix, "vote_")

	if _, exists := models.ActiveGame.Players[voterID]; !exists {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ Kamu sudah dieliminasi dan tidak bisa vote.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	if _, voted := playerVotes[userID]; voted {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ Kamu sudah memilih! Tidak bisa memilih lagi.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	playerVotes[userID] = voteTarget

	if voterID != userID {
		if voteTarget == string(models.ActiveGame.Players[userID].Symbol) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("🎃 <@%s>Kamu jujur, poinmu tetap -1!", userID),
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
		} else {
			models.ActiveGame.Players[voterID].Points+=2
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("🎭 <@%s>Kamu berbohong, poinmu +2!", userID),
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
		}
	} else {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("🎃 **<@%s>Kamu memilih _%s_, Good luck!**", userID, voteTarget),
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
	}
}
