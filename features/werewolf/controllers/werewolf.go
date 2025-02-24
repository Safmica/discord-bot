package controllers

import (
	"fmt"
	"strings"

	"github.com/Safmica/discord-bot/features/werewolf/models"
	"github.com/bwmarrin/discordgo"
)

var werewolfChannelID string
var werewolfID []string
var werewolfVotes = make(map[string]string)
var voteWerewolfCount = make(map[string]int)
var voteWerewolfMessageID string
var voteWerewolfStatus bool
var deathPlayerID []string

func StartWerewolfVoting(s *discordgo.Session, i *discordgo.InteractionCreate,channelID string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("🚨 Recovered from panic:", r)
		}
	}()
	voteLock.Lock()
	defer voteLock.Unlock()

	var voteOptions []discordgo.MessageComponent
	for id, player := range models.ActiveGame.Players {
		if player.Role != "werewolf" { 
			voteOptions = append(voteOptions, discordgo.Button{
				Label:    player.Username,
				Style:    discordgo.DangerButton,
				CustomID: fmt.Sprintf("vote_%s", id),
			})
		}
	}

	msg := &discordgo.Message{
		Content: "🔴 **Werewolf, pilih siapa yang akan kalian eliminasi!**",
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

func HandleWerewolfVote(s *discordgo.Session, i *discordgo.InteractionCreate, target string) {
	voteStatus = true
	voteLock.Lock()
	defer voteLock.Unlock()

	userID := i.Member.User.ID
	voteTarget := strings.TrimPrefix(target, "werewolf_vote_")

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

	found := false
	for _, id := range werewolfID {
		if id == i.Member.User.ID {
			found = true
			break
		}
	}

	if  !found {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ Kamu bukan Werewolf!",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}


	werewolfVotes[userID] = voteTarget
	voteWerewolfCount[voteTarget]++

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "✅ Vote kamu telah dicatat!",
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

	models.ActiveGame.Players[userID].VotingID = ""

	if len(playerVotes) == models.ActiveGame.Werewolf {
        voteWerewolfStatus = false
    }

	if !voteWerewolfStatus {
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

		if len(voteLeaders) == 1 {
			eliminatedPlayerID = voteLeaders[0]
			deathPlayerID = append(deathPlayerID, eliminatedPlayerID)
		}

		content := "🐺 **Sesi Werewolf Selesai!**"
		s.ChannelMessageEditComplex(&discordgo.MessageEdit{
			ID:         wwTextID,
			Channel:    models.ActiveGame.ID,
			Content:    &content,
			Components: &[]discordgo.MessageComponent{},
		})
	}
}