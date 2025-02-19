package bot

import (
	"fmt"
	"strings"

	"github.com/Safmica/discord-bot/config"
	"github.com/Safmica/discord-bot/features/undercover/controllers"
	"github.com/bwmarrin/discordgo"
)

var BotId string
var goBot *discordgo.Session

func Start() error{
	goBot, err := discordgo.New("Bot " + config.Config.Token)

	if err != nil {
		return err
	}

	goBot.Identify.Intents = discordgo.IntentsAll

	u, err := goBot.User("@me")

	if err != nil {
		return err
	}

	BotId = u.ID

	goBot.AddHandler(messageHandler)
	goBot.AddHandler(interactionHandler)


	err = goBot.Open()

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println("Bot is running!")

	return nil
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
    if m.Author.ID == BotId {
        return
    }

	if !strings.HasPrefix(m.Content, config.Config.BotPrefix) {
		return
	}

	command := strings.TrimPrefix(m.Content, config.Config.BotPrefix)

	switch command {
	case "start":
		controllers.StartGame(s, m, nil)
	}

	fmt.Println("Received message:", m.Content)

    if m.Content == "<@"+BotId+"> ping"{
        _, err := s.ChannelMessageSend(m.ChannelID, "pong!")
        if err != nil {
            fmt.Println("Error sending message:", err)
        }
    }

    if m.Content == config.Config.BotPrefix+"ping" {
        _, err := s.ChannelMessageSend(m.ChannelID, "pong!")
        if err != nil {
            fmt.Println("Error sending message:", err)
        }
    }
}

func interactionHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
    switch i.Type {
    case discordgo.InteractionMessageComponent:
        data := i.MessageComponentData()
        if data.CustomID == "join_game" {
            controllers.JoinGame(s, i)
        }
    }
}
