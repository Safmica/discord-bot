package controllers

import (
	"fmt"
	"strings"

	models "github.com/Safmica/discord-bot/features/jackheart/models"
	"github.com/bwmarrin/discordgo"
)

var jackVotes = make(map[string]string)
var voteCount = make(map[string]int)
var voteJackMessageID = ""
var playersVotes = 0

func JackVote(s *discordgo.Session, i *discordgo.InteractionCreate, prefix string) {
	voteStatus = true
	voteLock.Lock()
	defer voteLock.Unlock()

	userID := i.Member.User.ID
	voteTarget := strings.TrimPrefix(prefix, "jack_vote_")

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
	playersVotes++

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("✅ <@%s>Vote kamu telah dicatat!", userID),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})

	if playersVotes == len(models.ActiveGame.Players) {
		voteStatus = false
	}

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
				if eliminatedPlayerID == models.ActiveGame.Jackheart{
					gameStatus = false
					roles = "pawn"
				}
				delete(models.ActiveGame.Players, eliminatedPlayerID)
				eliminationMessage = fmt.Sprintf("☠️ <@%s> telah dieliminasi!", eliminatedPlayerID)
			} else {
				eliminationMessage = "🤷‍♂️ Pemain memilih skip! Tidak ada yang dieliminasi."
			}
		} else {
			eliminationMessage = "⚖️ Hasil voting seri! Tidak ada yang dieliminasi."
		}

		jackVotes = make(map[string]string)
		voteCount = make(map[string]int)
		voteJackMessageID = ""
		playersVotes = 0
		s.ChannelMessageSend(i.Interaction.ChannelID, eliminationMessage)
		SkipVotingCountdown()
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

		if voteJackMessageID != "" && voteStatus {
			_, err := s.ChannelMessageEdit(i.Interaction.ChannelID, voteJackMessageID, voteResults)
			if err != nil {
				fmt.Println("Gagal mengedit pesan voting:", err)
			}
		} else {
			msg, err := s.ChannelMessageSend(i.Interaction.ChannelID, voteResults)
			if err == nil {
				voteJackMessageID = msg.ID
			}
		}
	} 
}