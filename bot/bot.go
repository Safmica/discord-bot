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
	goBot.AddHandler(controllers.UndercoverHandler)


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
	command = strings.TrimSpace(command)
	
	switch {
	case command == "undercover":
		controllers.StartGame(s, m, nil)
	case strings.HasPrefix(command, "undercover config"):
		args := strings.TrimPrefix(command, "undercover config")
		args = strings.TrimSpace(args) // Hilangkan spasi ekstra sebelum mengirim ke controller
		controllers.ConfigUndercover(s, m ,args)
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