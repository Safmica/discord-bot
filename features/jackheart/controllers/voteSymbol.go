package controllers

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	models "github.com/Safmica/discord-bot/features/jackheart/models"
	"github.com/bwmarrin/discordgo"
)

var playerVotes = make(map[string]string)
var voteStatus bool
var voteLock sync.Mutex
var voteMessageID = ""
var showVotingID = ""
var Phase = "finish"

func startTurnBasedVoting(s *discordgo.Session, channelID string) {
	activeGame := models.ActiveGame
	players := activeGame.Players

	for _, player := range players {
		Phase = "ongoing"
		models.ActiveGame.NowPlaying = player.ID
		voteStatus = true

		content := fmt.Sprintf("📜 **Silahkan voting simbol <@%s> Dalam waktu 30 detik**\n _Selain <@%s>, voting bersifat opsional_", player.ID, player.ID)

		msg, err := s.ChannelMessageSendComplex(models.ActiveGame.ID, &discordgo.MessageSend{
			Content: content,
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label:    "Vote Symbol",
						Style:    discordgo.PrimaryButton,
						CustomID: "jackheart_view_vote_symbol",
					},
				}},
			},
		})
		if err != nil {
			fmt.Println("Gagal mengirim pesan:", err)
			return
		}
		models.ActiveGame.VotingID = msg.ID
		voteMessageID = models.ActiveGame.VotingID

		for i := 29; i >= 0; i-- {
			if !voteStatus {
				break
			}

			time.Sleep(1 * time.Second)

			voteResults := "📜 **Hasil Voting Sementara:**\n"
			for player, symbol := range playerVotes {
				voteResults += fmt.Sprintf("- <@%s> memilih **%s**\n", player, symbol)
			}

			content = fmt.Sprintf("📜 **Silahkan voting simbol <@%s> Dalam waktu %d detik**\n _Selain <@%s>, voting bersifat opsional_\n\n%s",
				player.ID, i, player.ID, voteResults)

			s.ChannelMessageEditComplex(&discordgo.MessageEdit{
				ID:      voteMessageID,
				Channel: channelID,
				Content: &content,
				Components: &[]discordgo.MessageComponent{
					discordgo.ActionsRow{Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "Vote Symbol",
							Style:    discordgo.PrimaryButton,
							CustomID: "jackheart_view_vote_symbol",
						},
					}},
				},
			})
		}
		content = fmt.Sprintf("🛑 **Voting untuk <@%s> telah berakhir!**", player.ID)
		s.ChannelMessageEditComplex(&discordgo.MessageEdit{
			ID:         voteMessageID,
			Channel:    channelID,
			Content:    &content,
			Components: &[]discordgo.MessageComponent{},
		})
		playerVotes = make(map[string]string)
	}
	Phase = "finish"

	for id, player := range models.ActiveGame.Players {
		rand.Shuffle(len(models.Symbols), func(i, j int) {
			models.Symbols[i], models.Symbols[j] = models.Symbols[j], models.Symbols[i]
		})
		player.Symbol = models.Symbols[0]
		models.ActiveGame.Players[id] = player
	}
	s.ChannelMessageSend(channelID, "✅ **Semua pemain telah menyelesaikan voting!**")
}

func ShowVote(s *discordgo.Session, i *discordgo.InteractionCreate) {
	voteLock.Lock()
	defer voteLock.Unlock()

	voterID := i.Member.User.ID
	userID := models.ActiveGame.NowPlaying
	players := models.ActiveGame.Players[userID]
	if showVotingID != "" {
		err := s.FollowupMessageDelete(i.Interaction, showVotingID)
		if err != nil {
			if strings.Contains(err.Error(), "Unknown Webhook") {
				fmt.Println("Pesan lama sudah tidak ada, lanjut tanpa menghapus.")
			} else {
				fmt.Println("Gagal menghapus pesan lama:", err)
			}
		}
		showVotingID = ""
		return
	}

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

	if _, voted := playerVotes[voterID]; voted {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ Kamu sudah memilih! Tidak bisa memilih lagi.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	models.ActiveGame.Players[voterID].Points--

	symbol := ""
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

	content := fmt.Sprintf("**Symbol milik <@%s> adalah %s**📜 **Pilihlah symbol untuknya**\n_pointmu berkurang -1_", players.ID, symbol)
	if players.ID == voterID {
		content = "📜 **Pilihlah symbol anda**\n"
	}

	buttons := []discordgo.MessageComponent{
		discordgo.Button{Label: "Heart ❤️", Style: discordgo.SecondaryButton, CustomID: "jackheart_vote_heart"},
		discordgo.Button{Label: "Spade ♠️", Style: discordgo.DangerButton, CustomID: "jackheart_vote_spade"},
		discordgo.Button{Label: "Diamond ♦️", Style: discordgo.DangerButton, CustomID: "jackheart_vote_diamond"},
		discordgo.Button{Label: "Club ♣️", Style: discordgo.DangerButton, CustomID: "jackheart_vote_club"},
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		fmt.Println("Gagal mengirim respon interaksi:", err)
		return
	}

	msg, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: content,
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{Components: buttons},
		},
		Flags: discordgo.MessageFlagsEphemeral,
	})

	if err != nil {
		fmt.Println("Gagal mengirim pesan:", err)
		return
	}

	if msg != nil {
		showVotingID = msg.ID
	}
}

func HandleVote(s *discordgo.Session, i *discordgo.InteractionCreate, prefix string) {
	fmt.Println(models.ActiveGame.Players)
	if models.ActiveGame == nil || !models.ActiveGame.Started {
		return
	}

	voteStatus = true
	voteLock.Lock()
	defer voteLock.Unlock()
	voterID := i.Member.User.ID
	userID := models.ActiveGame.NowPlaying

	voteTarget := strings.TrimPrefix(prefix, "jackheart_vote_")

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

	if _, voted := playerVotes[voterID]; voted {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ Kamu sudah memilih! Tidak bisa memilih lagi.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	playerVotes[voterID] = voteTarget

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
			models.ActiveGame.Players[voterID].Points += 2
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("🎭 <@%s>Kamu berbohong, poinmu +2!", userID),
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
		}
	} else {
		voteStatus = false
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("🎃 **<@%s>Kamu memilih _%s_, Good luck!**", userID, voteTarget),
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
	}
}
