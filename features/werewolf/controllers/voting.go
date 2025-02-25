package controllers

import (
	"fmt"
	"strings"

	"github.com/Safmica/discord-bot/features/undercover/models"
	"github.com/bwmarrin/discordgo"
)

var villagerVoteStatus = false

func HandleVote(s *discordgo.Session, i *discordgo.InteractionCreate, prefix string) {
	voteStatus = true
	villagerVoteStatus = false
	voteLock.Lock()
	defer voteLock.Unlock()

	userID := i.Member.User.ID
	voteTarget := strings.TrimPrefix(prefix, "werewolf_vote_")

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

	playerVotes[userID] = voteTarget
	voteCount[voteTarget]++

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("✅ <@%s>Vote kamu telah dicatat!", userID),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})

	gameEnded := false
	if len(playerVotes) == len(models.ActiveGame.Players) {
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
				eliminationMessage = fmt.Sprintf("☠️ **<@%s> telah dieliminasi!**", eliminatedPlayerID)
				delete(models.ActiveGame.Players, eliminatedPlayerID)
			} else {
				eliminationMessage = "🤷‍♂️ Pemain memilih skip! Tidak ada yang dieliminasi."
			}
		} else {
			eliminationMessage = "⚖️ Hasil voting seri! Tidak ada yang dieliminasi."
		}

		villagerCount, werewolfCount := 0, 0
		for _, player := range models.ActiveGame.Players {
			if player.Role == "werewolf" {
				werewolfCount++
			} else {
				villagerCount++
			}
		}

		var endMessage string
		if werewolfCount == 0 {
			endMessage = "🎉 **Villager menang!** Semua Werewolf telah dieliminasi."
			gameEnded = true
		} else if werewolfCount >= villagerCount {
			endMessage = "🐺 **Werewolf menang!** Mereka berhasil menguasai desa."
			gameEnded = true
		}

		if gameEnded {
			models.ActiveGame = nil
			s.ChannelMessageSend(i.Interaction.ChannelID, endMessage)
			villagerVoteStatus = true
		} else {
			s.ChannelMessageSend(i.Interaction.ChannelID, eliminationMessage)
			villagerVoteStatus = true
		}
	}

	if !gameEnded {
		var voteResults string
		voteResults = "📊 **Hasil Voting Sementara:**\n"
		for playerID, count := range voteCount {
			if playerID == "skip" {
				voteResults += fmt.Sprintf("- %s: %d suara\n", playerID, count)
			} else {
				voteResults += fmt.Sprintf("- <@%s>: %d suara\n", playerID, count)
			}
		}

		unvotedPlayers := []string{}
		for playerID := range models.ActiveGame.Players {
			if _, voted := playerVotes[playerID]; !voted {
				unvotedPlayers = append(unvotedPlayers, fmt.Sprintf("<@%s>", playerID))
			}
		}

		if len(unvotedPlayers) > 0 {
			voteResults += "\n⏳ **Pemain yang belum vote:** " + strings.Join(unvotedPlayers, ", ")
		} else {
			voteResults += "\n✅ Semua pemain sudah vote!"
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
			playerVotes = make(map[string]string)
			voteCount = make(map[string]int)
			voteMessageID = ""
		} else {
			msg, err := s.ChannelMessageSend(i.Interaction.ChannelID, voteResults)
			if err == nil {
				voteMessageID = msg.ID
			}
		}
	} else {
		voteResults := "📊 **Game Berakhir**\n"
		villagerVoteStatus = true
		playerVotes = make(map[string]string)
		voteCount = make(map[string]int)
		voteMessageID = ""
		gameStatus = false
		villagerVoteStatus = false

		msg := &discordgo.MessageSend{
			Content: voteResults,
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "Play Again",
							Style:    discordgo.PrimaryButton,
							CustomID: "play_again_werewolf",
						},
					},
				},
			},
		}

		s.ChannelMessageSendComplex(i.Interaction.ChannelID, &discordgo.MessageSend{
			Content:    msg.Content,
			Components: msg.Components,
		})
	}
}