package controllers

import (
	"fmt"
	"strings"

	models "github.com/Safmica/discord-bot/features/jackheart/models"
	"github.com/bwmarrin/discordgo"
)

var jackVotes = make(map[string]string)
var voteCount = make(map[string]int)

func JackVote(s *discordgo.Session, i *discordgo.InteractionCreate, prefix string) {
	voteStatus = true
	voteLock.Lock()
	defer voteLock.Unlock()

	userID := i.Member.User.ID
	voteTarget := strings.TrimPrefix(prefix, "vote_")

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

	if _, voted := jackVotes[userID]; voted {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ Kamu sudah memilih! Tidak bisa memilih lagi.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	jackVotes[userID] = voteTarget
	voteCount[voteTarget]++

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("✅ <@%s>Vote kamu telah dicatat!", userID),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})

	if !voteStatus {
		maxVotes := 0
		voteLeaders := []string{}

		for playerID, count := range voteCount {
			if count > maxVotes {
				maxVotes = count
				voteLeaders = []string{playerID}
			} else if count == maxVotes {
				voteLeaders = append(voteLeaders, playerID)
			}
		}

		var eliminatedPlayerID string
		eliminationMessage := ""

		if len(voteLeaders) == 1 {
			eliminatedPlayerID = voteLeaders[0]
			if eliminatedPlayerID != "skip" {
				delete(models.ActiveGame.Players, eliminatedPlayerID)
				eliminationMessage = fmt.Sprintf("☠️ <@%s> telah dieliminasi!", eliminatedPlayerID)
			} else {
				eliminationMessage = "🤷‍♂️ Pemain memilih skip! Tidak ada yang dieliminasi."
			}
		} else {
			eliminationMessage = "⚖️ Hasil voting seri! Tidak ada yang dieliminasi."
		}

		s.ChannelMessageSend(i.Interaction.ChannelID, eliminationMessage)
		SendVotingMessage(s, i, i.Interaction.ChannelID)
	} else {
		var voteResults string
		voteResults = "📊 **Hasil Voting Sementara:**\n"
		for playerID, count := range voteCount {
			if playerID == "skip" {
				voteResults += fmt.Sprintf("- %s: %d suara\n", playerID, count)
			} else {
				voteResults += fmt.Sprintf("- <@%s>: %d suara\n", playerID, count)
			}
		}

		if voteMessageID != "" && voteStatus {
			_, err := s.ChannelMessageEdit(i.Interaction.ChannelID, voteMessageID, voteResults)
			if err != nil {
				fmt.Println("Gagal mengedit pesan voting:", err)
			}
		} else if voteMessageID != "" && !voteStatus {
			voteResults = "📊 **Hasil Voting Akhir:**\n"
			for playerID, count := range voteCount {
				if playerID == "skip" {
					voteResults += fmt.Sprintf("- %s: %d suara\n", playerID, count)
				} else {
					voteResults += fmt.Sprintf("- <@%s>: %d suara\n", playerID, count)
				}
			}
			_, err := s.ChannelMessageEdit(i.Interaction.ChannelID, voteMessageID, voteResults)
			if err != nil {
				fmt.Println("Gagal mengedit pesan voting:", err)
			}
			jackVotes = make(map[string]string)
			voteCount = make(map[string]int)
			voteMessageID = ""
		} else {
			msg, err := s.ChannelMessageSend(i.Interaction.ChannelID, voteResults)
			if err == nil {
				voteMessageID = msg.ID
			}
		}
	} 
}