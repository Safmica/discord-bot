package controllers

import (
	"fmt"
	"strings"
	"sync"
	"time"

	models "github.com/Safmica/discord-bot/features/jackheart/models"
	"github.com/bwmarrin/discordgo"
)

var playerVotes = make(map[string]string)
var voteStatus bool
var voteLock sync.Mutex

func startTurnBasedVoting(s *discordgo.Session, channelID string) {
	activeGame := models.ActiveGame
	players := activeGame.Players

	for _, player := range players {
		var voteMessageID = models.ActiveGame.VotingID
		models.ActiveGame.NowPlaying = player.ID
		voteStatus = true
	
		content := fmt.Sprintf("📜 **Silahkan voting simbol <@%s> Dalam waktu 30 detik**\n _Selain <@%s>, voting bersifat opsional_", player.ID, player.ID)
	
		// Kirim pesan pertama kali
		if voteMessageID == "" {
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
			content = fmt.Sprintf("📜 **Silahkan voting simbol <@%s> Dalam waktu %d detik**\n _Selain <@%s>, voting bersifat opsional_\n\n%s", 
			player.ID,i, player.ID, voteResults)
	
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
		playerVotes = make(map[string]string)
	}
	
	// for id, player := range models.ActiveGame.Players {
	// 	rand.Shuffle(len(models.Symbols), func(i, j int) { 
	// 		models.Symbols[i], models.Symbols[j] = models.Symbols[j], models.Symbols[i] 
	// 	})
	// 	player.Symbol = models.Symbols[0] 
	// 	models.ActiveGame.Players[id] = player
	// }	
	s.ChannelMessageSend(channelID, "✅ **Semua pemain telah menyelesaikan voting!**")
}	

func ShowVote(s *discordgo.Session, i *discordgo.InteractionCreate) {
	voteLock.Lock()
	defer voteLock.Unlock()

	voterID := i.Member.User.ID
	models.ActiveGame.Players[voterID].Points--
	userID := models.ActiveGame.NowPlaying

	players := models.ActiveGame.Players[userID]
	symbol := ""
	content := ""
	fmt.Println(players.Symbol)
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
	if players.ID == voterID {
		content = "📜 **Pilihlah symbol anda**\n"
	} else {
		content = fmt.Sprintf("**Symbol milik <@%s> adalah %s**📜 **Pilihlah symbol untuknya**\n_pointmu berkurang -1_", players.ID, symbol)
	}

	var buttons []discordgo.MessageComponent

	buttons = append(buttons, discordgo.Button{
		Label:    "Heart ❤️",
		Style:    discordgo.SecondaryButton,
		CustomID: "jackheart_vote_heart",
	})

	buttons = append(buttons, discordgo.Button{
		Label:    "Spade ♠️",
		Style:    discordgo.DangerButton,
		CustomID: "jackheart_vote_spade",
	})

	buttons = append(buttons, discordgo.Button{
		Label:    "Diamond ♦️",
		Style:    discordgo.DangerButton,
		CustomID: "jackheart_vote_diamond",
	})

	buttons = append(buttons, discordgo.Button{
		Label:    "Club ♣️",
		Style:    discordgo.DangerButton,
		CustomID: "jackheart_vote_club",
	})

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
	

	msg, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: content,
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{Components: buttons},
		},
		Flags:   discordgo.MessageFlagsEphemeral,
	})

	if err != nil {
		fmt.Println("Gagal mengirim pesan:", err)
	}

	lastVoteMessageID = msg.ID
	go func() {
		timeout := 30 // Batas waktu 30 detik
		for timeout > 0 {
			time.Sleep(1 * time.Second)
			if !voteStatus { 
				break
			}
			timeout--
		}
	
		emptyString := "_⛔Vote selesai_"
		_, err := s.FollowupMessageEdit(i.Interaction, lastVoteMessageID, &discordgo.WebhookEdit{
			Content: &emptyString, // Mengosongkan pesan
			Components: &[]discordgo.MessageComponent{}, // Menghapus tombol
		})
		if err != nil {
			fmt.Println("Gagal mengedit pesan:", err)
		}
	}()	
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
