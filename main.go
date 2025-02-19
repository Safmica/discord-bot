package main

import (
	"fmt"
	"log"

	"github.com/Safmica/discord-bot/bot"
	"github.com/Safmica/discord-bot/config"
)

func main() {
	if err := config.ReadConfig(); err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	if err := bot.Start(); err != nil {
		log.Fatalf("Failed to start bot: %v", err)
	}

	fmt.Println("Bot is running... Press Ctrl+C to exit.")

	select {}
}
