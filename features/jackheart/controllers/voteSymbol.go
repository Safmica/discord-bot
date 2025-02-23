package controllers

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"sync"
	"time"

	models "github.com/Safmica/discord-bot/features/jackheart/models"
	"github.com/bwmarrin/discordgo"
)

var playerVotes = make(map[string]string)
var voteStatus bool
var voteLock sync.Mutex
var Phase = "finish"
var voteMessageID = ""
var gameStatus = true
var skipTimer = make(chan struct{})
var phaseNumber = 1
var roles = "jackheart"

func startTurnBasedVoting(s *discordgo.Session, channelID string) {
	activeGame := models.ActiveGame
	players := activeGame.Players

	content := ""
	for gameStatus {
		for _, players := range models.ActiveGame.Players {
			rand.Shuffle(len(models.Symbols), func(i, j int) { models.Symbols[i], models.Symbols[j] = models.Symbols[j], models.Symbols[i] })
			players.Symbol = models.Symbols[0]
		}

		if models.ActiveGame.TempVoteVotingID != "" {
			content = "🛑 **Voting untuk telah berakhir!**"
			s.ChannelMessageEditComplex(&discordgo.MessageEdit{
				ID:         models.ActiveGame.JackVoteID,
				Channel:    channelID,
				Content:    &content,
				Components: &[]discordgo.MessageComponent{},
			})
		}
		if phaseNumber > 1 {
			msg, _ := s.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
				Content: fmt.Sprintf("**MEMASUKI FASE -%d**", phaseNumber),
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "View Dashboard",
							Style:    discordgo.PrimaryButton,
							CustomID: "view_dashboard_jackheart",
						},
					}}}})
			models.ActiveGame.AnnounceID = msg.ID
		} else {
			s.ChannelMessageSend(channelID, fmt.Sprintf("**MEMASUKI FASE -%d**", phaseNumber))
		}
		phaseNumber++
		for _, player := range players {
			Phase = "ongoing"
			models.ActiveGame.NowPlaying = player.ID
			voteStatus = true

			content = fmt.Sprintf("📜 **Silahkan voting simbol <@%s> Dalam waktu 100 detik**\n _Selain <@%s>, voting bersifat opsional_\n\n📜 **Saran dari player lain:** (_hati hati mereka mungkin berbohong_)\n", player.ID, player.ID)

			msg, err := s.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
				Content: content,
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "Vote Symbol",
							Style:    discordgo.PrimaryButton,
							CustomID: "jackheart_view_vote_symbol",
						},
					}}}})

			if err != nil {
				fmt.Println("Gagal mengirim pesan:", err)
				return
			}
			models.ActiveGame.VotingID = msg.ID
			voteMessageID = models.ActiveGame.VotingID

			for i := 99; i >= 0; i-- {
				voteLock.Lock()
				if !voteStatus {
					voteLock.Unlock()
					break
				}
				voteLock.Unlock()

				time.Sleep(1 * time.Second)

				players := make([]string, 0, len(playerVotes))
				for player := range playerVotes {
					players = append(players, player)
				}

				sort.Strings(players)

				voteResults := "📜 **Saran dari player lain:** (_hati hati mereka mungkin berbohong_)\n"
				for _, player := range players {
					symbol := ""
					switch playerVotes[player] {
					case "heart":
						symbol = "**Heart ❤️**"
					case "diamond":
						symbol = "**Diamond ♦️**"
					case "club":
						symbol = "**Club ♣️**"
					case "spade":
						symbol = "**Spade ♠️**"
					}
					voteResults += fmt.Sprintf("- Simbolmu adalah **%s** kata <@%s>\n", symbol, player)
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
						}}}})
			}
			if models.ActiveGame.TempVoteVotingID != ""{
				content = "🛑 **Voting untuk telah berakhir!**"
				s.ChannelMessageEditComplex(&discordgo.MessageEdit{
					ID:         models.ActiveGame.TempVoteVotingID,
					Channel:    channelID,
					Content:    &content,
					Components: &[]discordgo.MessageComponent{},
				})
			}
			models.ActiveGame.TempVoteVotingID = msg.ID

			content = fmt.Sprintf("🛑 **Voting untuk <@%s> telah berakhir!**\n _Tekan **View Dashboard** jika ingin melihat point anda_", player.ID)
			s.ChannelMessageEditComplex(&discordgo.MessageEdit{
				ID:      voteMessageID,
				Channel: channelID,
				Content: &content,
				Components: &[]discordgo.MessageComponent{
					discordgo.ActionsRow{Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "View Dashboard",
							Style:    discordgo.PrimaryButton,
							CustomID: "view_dashboard_jackheart",
						},
					}}}})
		}

		playerVotes = make(map[string]string)

		Phase = "finish"

		message := "✅ **Semua pemain telah menyelesaikan voting!\nHASIL SEMENTARA**(_Diantara kalian terdapat jack! Vote Pemain mencurigakan!!_)\n"

		for id, player := range models.ActiveGame.Players {
			if player.Symbol == player.SymbolVoted && player.Points >= 0 {
				message += fmt.Sprintf("- <@%s> 🟢 **Hidup** | **Poin:** %d\n", id, player.Points)
			} else {
				message += fmt.Sprintf("- <@%s> ☠️ **Mati**\n", id)
				delete(models.ActiveGame.Players, id)
			}
		}

		jackHeart := false
		for _, player := range models.ActiveGame.Players {
			player.Symbol = ""
			player.SymbolVoted = ""
			if player.Points >= models.ActiveGame.MaxPoints {
				gameStatus = false
				roles = string(player.Role)
			}
			if player.ID == models.ActiveGame.Jackheart {
				jackHeart = true
			}
		}

		if !jackHeart {
			gameStatus = false
			roles = "pawn"
		}
		if jackHeart {
			if len(models.ActiveGame.Players) == 1 {
				gameStatus = false
				roles = "jackheart"
			}
		}

		if len(models.ActiveGame.Players) <= 0 {
			gameStatus = false
			roles = "draw"
		}

		if gameStatus {
			var components []discordgo.MessageComponent
			var buttons []discordgo.MessageComponent

			for _, p := range players {

				buttons = append(buttons, discordgo.Button{
					Label:    p.Username,
					Style:    discordgo.PrimaryButton,
					CustomID: "jack_vote_" + p.ID,
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
				CustomID: "jack_vote_skip",
			})

			if len(buttons) == 5 {
				components = append(components, discordgo.ActionsRow{Components: buttons})
				buttons = []discordgo.MessageComponent{}
			}

			if len(buttons) > 0 {
				components = append(components, discordgo.ActionsRow{Components: buttons})
			}

			msg , err := s.ChannelMessageSendComplex(models.ActiveGame.ID, &discordgo.MessageSend{
				Content:    message,
				Components: components,
			})

			models.ActiveGame.JackVoteID = msg.ID

			if err != nil {
				fmt.Println("Gagal mengirim pesan:", err)
			}

			if !gameStatus {
				goto EndGame
			}

			StartVotingCountdown(s)
			time.Sleep(1 * time.Second)
			if votedPlayer != "" {
				if models.ActiveGame.Players[votedPlayer].Points < 0 {
					if votedPlayer == models.ActiveGame.Jackheart {
						gameStatus = false
						roles = "pawn"
					}
					message = fmt.Sprintf("- <@%s> ☠️ **Mati**\n", votedPlayer)
					delete(models.ActiveGame.Players, votedPlayer)
					votedPlayer = ""
					s.ChannelMessageSend(channelID, message)
				}
			}
			content := "🛑 **Voting untuk telah berakhir!**"
			s.ChannelMessageEditComplex(&discordgo.MessageEdit{
				ID:         models.ActiveGame.TempVoteVotingID,
				Channel:    channelID,
				Content:    &content,
				Components: &[]discordgo.MessageComponent{},
			})
		}

	EndGame:
		if !gameStatus {
			content := "🛑 **Voting untuk telah berakhir!**"
			s.ChannelMessageEditComplex(&discordgo.MessageEdit{
				ID:         models.ActiveGame.TempVoteVotingID,
				Channel:    channelID,
				Content:    &content,
				Components: &[]discordgo.MessageComponent{},
			})
			message := ""
			if roles == "pawn" {
				message = fmt.Sprintf("**🎊Pion Menang! Jackheart adalah <@%s>🎉**", models.ActiveGame.Jackheart)
			} else if roles == "jackheart" {
				message = fmt.Sprintf("**🎊Jackheart Menang! Jackheart adalah <@%s>🎉**", models.ActiveGame.Jackheart)
			} else {
				message = "**🤣HAHAHA Kalian semuat mati, tidak ada yang menang🤣**"
			}
			msg := &discordgo.MessageSend{
				Content: message,
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.Button{
								Label:    "Play Again",
								Style:    discordgo.PrimaryButton,
								CustomID: "play_again_jackheart",
							},
						},
					},
				},
			}
			models.ActiveGame = nil
			phaseNumber = 1
			Phase = "finish"
			voteMessageID = ""
			gameStatus = false
			playerVotes = make(map[string]string)
			playerReady = 0
			jackVotes = make(map[string]string)
			voteCount = make(map[string]int)
			voteJackMessageID = ""
			playersVotes = 0
			s.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
				Content:    msg.Content,
				Components: msg.Components,
			})
		}
	}
}

func ShowVote(s *discordgo.Session, i *discordgo.InteractionCreate) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("🚨 Recovered from panic:", r)
		}
	}()
	voteLock.Lock()
	defer voteLock.Unlock()

	voterID := i.Member.User.ID
	userID := models.ActiveGame.NowPlaying
	players := models.ActiveGame.Players[userID]
	voter := models.ActiveGame.Players[voterID]
	if models.ActiveGame.Players[voterID].VoteID != "" {
		err := s.FollowupMessageDelete(i.Interaction, models.ActiveGame.Players[voterID].VoteID)
		if err != nil {
			if strings.Contains(err.Error(), "Unknown Webhook") {
				fmt.Println("Pesan lama sudah tidak ada, lanjut tanpa menghapus.")
			} else {
				fmt.Println("Gagal menghapus pesan lama:", err)
			}
		}
		models.ActiveGame.Players[voterID].VoteID = ""
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

	if voter.ID != players.ID {
		models.ActiveGame.Players[voterID].Points--
	}

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
		discordgo.Button{Label: "Spade ♠️", Style: discordgo.SecondaryButton, CustomID: "jackheart_vote_spade"},
		discordgo.Button{Label: "Diamond ♦️", Style: discordgo.SecondaryButton, CustomID: "jackheart_vote_diamond"},
		discordgo.Button{Label: "Club ♣️", Style: discordgo.SecondaryButton, CustomID: "jackheart_vote_club"},
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

	models.ActiveGame.Players[voterID].VoteID = msg.ID
	models.ActiveGame.Players[voterID].VotingPlayer = userID
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

	if models.ActiveGame.Players[voterID].VotingPlayer != models.ActiveGame.NowPlaying {
		if models.ActiveGame.Players[voterID].VotingPlayer != userID {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "❌ Kamu tidak bisa memilih disini.",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			return
		}
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
					Content: "🎃 Kamu jujur, poinmu tetap -1!",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
		} else {
			models.ActiveGame.Players[voterID].Points += 2
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "🎭 Kamu berbohong, poinmu +2!",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
		}
	} else {
		voteStatus = false
		switch voteTarget {
		case "heart":
			models.ActiveGame.Players[userID].SymbolVoted = models.Heart
		case "spade":
			models.ActiveGame.Players[userID].SymbolVoted = models.Spade
		case "club":
			models.ActiveGame.Players[userID].SymbolVoted = models.Club
		case "diamond":
			models.ActiveGame.Players[userID].SymbolVoted = models.Diamond
		}

		playerVotes = make(map[string]string)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "🎃 **Pilihanmu Berhasil Dicatat!**",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
	}
}

func StartVotingCountdown(s *discordgo.Session) {
	timer := time.NewTimer(300 * time.Second)

	select {
	case <-timer.C:
		fmt.Println("Waktu habis, lanjut ke tahap berikutnya!")
	case <-skipTimer:
		fmt.Println("Timer di-skip, lanjut ke tahap berikutnya lebih cepat!")
		timer.Stop()
	}
}

func SkipVotingCountdown() {
	select {
	case skipTimer <- struct{}{}:
	default:
	}
}
