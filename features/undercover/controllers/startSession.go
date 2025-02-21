package controllers

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"

	models "github.com/Safmica/discord-bot/features/undercover"
	"github.com/bwmarrin/discordgo"
)

func StartGameSession(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if models.ActiveGame == nil || models.ActiveGame.Started {
		return
	}

	playerCount := len(models.ActiveGame.Players)
	if playerCount < 4 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ Minimal 4 pemain diperlukan untuk memulai game.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	models.ActiveGame.Started = true

	players := make([]*models.Player, 0, playerCount)
	for _, p := range models.ActiveGame.Players {
		players = append(players, p)
	}

	rand.Shuffle(len(players), func(i, j int) { players[i], players[j] = players[j], players[i] })

	civilianWord := "Apple"
	undercoverWord := "Orange"

	for i := 0; i < models.ActiveGame.Undercover; i++ {
		players[i].Role = models.Undercover
	}

    for i := models.ActiveGame.Undercover; i < len(players); i++ {
		players[i].Role = models.Civilian
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "🚀 Game telah dimulai! Peranmu akan dikirim secara rahasia melalui DM.",
		},
	})

	for _, p := range players {
		word := ""
		switch p.Role {
		case models.Civilian:
			word = civilianWord
		case models.Undercover:
			word = undercoverWord
		}

		dmChannel, err := s.UserChannelCreate(p.ID)
		if err != nil {
			fmt.Println("Gagal membuat DM channel:", err)
			continue
		}

		_, err = s.ChannelMessageSend(dmChannel.ID, fmt.Sprintf("🔐 **Kata Rahasia Kamu:** %s", word))
		if err != nil {
			fmt.Println("Gagal mengirim DM:", err)
		}
	}

    playerList := "📜 **Daftar Pemain:**\n"
    var components []discordgo.MessageComponent
    var buttons []discordgo.MessageComponent

    for _, p := range players {
        playerList += fmt.Sprintf("- <@%s>\n", p.ID)
    
        buttons = append(buttons, discordgo.Button{
            Label:    p.Username,
            Style:    discordgo.PrimaryButton,
            CustomID: "vote_" + p.ID,
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
        CustomID: "vote_skip",
    })

    if len(buttons) == 5 {
        components = append(components, discordgo.ActionsRow{Components: buttons})
        buttons = []discordgo.MessageComponent{}
    }

    buttons = append(buttons, discordgo.Button{
        Label:    "Close Vote",
        Style:    discordgo.DangerButton,
        CustomID: "vote_close",
    })

    if len(buttons) > 0 {
        components = append(components, discordgo.ActionsRow{Components: buttons})
    }

    msg, err := s.ChannelMessageSendComplex(i.Interaction.ChannelID, &discordgo.MessageSend{
        Content: playerList + "\n🗳 **Silakan pilih aksi berikut!**",
        Components: components,
    })
    
    if err != nil {
        fmt.Println("Gagal mengirim pesan:", err)
    }

	lastVoteMessageID = msg.ID
}

var playerVotes = make(map[string]string)
var voteCount = make(map[string]int)
var voteMessageID string
var voteStatus bool
var voteLock sync.Mutex

func HandleVote(s *discordgo.Session, i *discordgo.InteractionCreate, prefix string) {
    voteStatus = true
    voteLock.Lock()
    defer voteLock.Unlock()

    userID := i.Member.User.ID
    voteTarget := strings.TrimPrefix(prefix, "vote_")

    if voteTarget == "close" && userID != models.ActiveGame.HostID {
        s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
            Type: discordgo.InteractionResponseChannelMessageWithSource,
            Data: &discordgo.InteractionResponseData{
                Content: "❌ Hanya host yang bisa melakukan close vote.",
                Flags:   discordgo.MessageFlagsEphemeral,
            },
        })
        return
    }

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

    if voteTarget != "close" {
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
    }

    playerVotes[userID] = voteTarget
    voteCount[voteTarget]++
    if voteTarget == "close" {
        voteStatus = false
    }

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

        if  len(voteLeaders) == 1  {
            eliminatedPlayerID = voteLeaders[0]
            if eliminatedPlayerID != "skip" {
                if models.ActiveGame.ShowRoles {
                    elimnatedPlayer := models.ActiveGame.Players[eliminatedPlayerID]
                    delete(models.ActiveGame.Players, eliminatedPlayerID)
                    eliminationMessage = fmt.Sprintf("☠️ **<@%s> telah dieliminasi! Dan dia merupakan _%s_**", eliminatedPlayerID, elimnatedPlayer.Role)
                } else {
                    delete(models.ActiveGame.Players, eliminatedPlayerID)
                    eliminationMessage = fmt.Sprintf("☠️ <@%s> telah dieliminasi!", eliminatedPlayerID)
                }
            } else {
                eliminationMessage = "🤷‍♂️ Pemain memilih skip! Tidak ada yang dieliminasi."
            }
        } else {
            eliminationMessage = "⚖️ Hasil voting seri! Tidak ada yang dieliminasi."
        }

        civilianCount, undercoverCount := 0, 0
        for _, player := range models.ActiveGame.Players {
            if player.Role == models.Civilian {
                civilianCount++
            } else if player.Role == models.Undercover {
                undercoverCount++
            }
        }

        var endMessage string
        if undercoverCount == 0 {
            endMessage = "🎉 **Civilian menang!** Semua Undercover telah dieliminasi."
            gameEnded = true
        } else if undercoverCount >= civilianCount {
            endMessage = "🤫 **Undercover menang!** Mereka berhasil menguasai permainan."
            gameEnded = true
        }

        if gameEnded {
            models.ActiveGame = nil
            s.ChannelMessageSend(i.Interaction.ChannelID, endMessage)
        } else {
            s.ChannelMessageSend(i.Interaction.ChannelID, eliminationMessage)
            SendVotingMessage(s, i, i.Interaction.ChannelID)
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

        if voteMessageID != "" && voteStatus{
            _, err := s.ChannelMessageEdit(i.Interaction.ChannelID, voteMessageID, voteResults)
            if err != nil {
                fmt.Println("Gagal mengedit pesan voting:", err)
            }
        } else if voteMessageID != "" && !voteStatus{
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
        }else {
            msg, err := s.ChannelMessageSend(i.Interaction.ChannelID, voteResults)
            if err == nil {
                voteMessageID = msg.ID
            }
        }
    } else {
        voteResults := "📊 **Game Berakhir**\n"
        playerVotes = make(map[string]string)
        voteCount = make(map[string]int)
        voteMessageID = ""

        msg := &discordgo.MessageSend{
            Content: voteResults,
            Components: []discordgo.MessageComponent{
                discordgo.ActionsRow{
                    Components: []discordgo.MessageComponent{
                        discordgo.Button{
                            Label:    "Play Again",
                            Style:    discordgo.PrimaryButton,
                            CustomID: "play_again",
                        },
                    },
                },
            },
        }

        s.ChannelMessageSendComplex(i.Interaction.ChannelID, &discordgo.MessageSend{
            Content: msg.Content,
            Components: msg.Components,
        })
    }
}

var lastVoteMessageID string

func SendVotingMessage(s *discordgo.Session,i *discordgo.InteractionCreate, channelID string) {
    if models.ActiveGame == nil || !models.ActiveGame.Started {
        return
    }

    CloseVoting(s,i, channelID)

    playerList := "📜 **Daftar Pemain yang Masih Hidup:**\n"
    var components []discordgo.MessageComponent
    var buttons []discordgo.MessageComponent
    
    for _, p := range models.ActiveGame.Players {
        playerList += fmt.Sprintf("- <@%s>\n", p.ID)
    
        buttons = append(buttons, discordgo.Button{
            Label:    p.Username,
            Style:    discordgo.PrimaryButton,
            CustomID: "vote_" + p.ID,
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
        CustomID: "vote_skip",
    })

    if len(buttons) == 5 {
        components = append(components, discordgo.ActionsRow{Components: buttons})
        buttons = []discordgo.MessageComponent{}
    }

    buttons = append(buttons, discordgo.Button{
        Label:    "Close Vote",
        Style:    discordgo.DangerButton,
        CustomID: "vote_close",
    })

    if len(buttons) > 0 {
        components = append(components, discordgo.ActionsRow{Components: buttons})
    }

    s.ChannelMessageSendComplex(i.Interaction.ChannelID, &discordgo.MessageSend{
        Content: playerList + "\n🗳 **Silakan pilih aksi berikut!**",
        Components: components,
    })
    
}

func CloseVoting(s *discordgo.Session,i *discordgo.InteractionCreate, channelID string) {
	if lastVoteMessageID != "" {
		content := "🗳 **Voting telah selesai!**"
		_, err := s.ChannelMessageEditComplex(&discordgo.MessageEdit{
			ID:      lastVoteMessageID,
			Channel: i.Interaction.ChannelID,
			Content: &content, 
			Components: &[]discordgo.MessageComponent{},  
		})
		if err != nil {
			fmt.Println("Gagal menghapus tombol dari pesan voting:", err)
		}
		lastVoteMessageID = "" 
	}
}
