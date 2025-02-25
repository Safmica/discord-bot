package controllers

import (
	"fmt"
	"strings"

	"github.com/Safmica/discord-bot/features/werewolf/models"
	"github.com/bwmarrin/discordgo"
)

func StartSeerVoting(s *discordgo.Session, i *discordgo.InteractionCreate,channelID string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("🚨 Recovered from panic:", r)
		}
	}()
	voteLock.Lock()
	defer voteLock.Unlock()

	userID := i.Member.User.ID

	if  userID != models.ActiveGame.SeerID {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ Kamu bukan Seer!",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	if _, exists := models.ActiveGame.Players[userID]; !exists {
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

	for _, id := range deathPlayerID {
		if id == userID {
			msg := &discordgo.Message{
				Content: "🐺 **Kamu telah dimangsa werewolf!**",
				Components: []discordgo.MessageComponent{},
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
		
			_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
				Content: msg.Content,
				Components: msg.Components,
				Flags: discordgo.MessageFlagsEphemeral,
			})
		
			if err != nil {
				fmt.Println("Gagal mengirim pesan:", err)
				return
			}
			content := "🧙🏻‍♀️ **Sesi Seer Selesai!**"
			s.ChannelMessageEditComplex(&discordgo.MessageEdit{
				ID:         seerTextID,
				Channel:    models.ActiveGame.ID,
				Content:    &content,
				Components: &[]discordgo.MessageComponent{},
			})
		
			seerVoteStatus = true
			return
		}
	}

	var voteOptions []discordgo.MessageComponent
	for id, player := range models.ActiveGame.Players {
		if player.ID != userID { 
			voteOptions = append(voteOptions, discordgo.Button{
				Label:    player.Username,
				Style:    discordgo.DangerButton,
				CustomID: fmt.Sprintf("seer_vote_%s", id),
			})
		}
	}

	msg := &discordgo.Message{
		Content: "🔮 **Seer, pilih siapa yang akan kamu ramal!**",
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{Components: voteOptions},
		},
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

	msg, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: msg.Content,
		Components: msg.Components,
		Flags: discordgo.MessageFlagsEphemeral,
	})

	if err != nil {
		fmt.Println("Gagal mengirim pesan:", err)
		return
	}

	models.ActiveGame.Players[i.Member.User.ID].VotingID = msg.ID
}

func HandleSeerVote(s *discordgo.Session, i *discordgo.InteractionCreate, target string) {
	voteLock.Lock()
	defer voteLock.Unlock()

	userID := i.Member.User.ID
	voteTarget := strings.TrimPrefix(target, "seer_vote_")

	if models.ActiveGame == nil || !models.ActiveGame.Started {
		return
	}

	if _, exists := models.ActiveGame.Players[userID]; !exists {
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

	if  userID != models.ActiveGame.SeerID {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ Kamu bukan Seer!",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	selected := models.ActiveGame.Players[voteTarget]

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Role milik %s adalah **%s**!", selected.Username, selected.Role),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})

	err := s.FollowupMessageDelete(i.Interaction, models.ActiveGame.Players[userID].VotingID)
	if err != nil {
		if strings.Contains(err.Error(), "Unknown Webhook") {
			fmt.Println("Pesan lama sudah tidak ada, lanjut tanpa menghapus.")
		} else {
			fmt.Println("Gagal menghapus pesan lama:", err)
		}
	}

	content := "🧙🏻‍♀️ **Sesi Seer Selesai!**"
	s.ChannelMessageEditComplex(&discordgo.MessageEdit{
		ID:         seerTextID,
		Channel:    models.ActiveGame.ID,
		Content:    &content,
		Components: &[]discordgo.MessageComponent{},
	})

	seerVoteStatus = true
}