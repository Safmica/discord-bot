package config

import (
	"os"

	"github.com/Safmica/discord-bot/models"
	"github.com/joho/godotenv"
)

var Config = &models.Config{}

func ReadConfig() error {
	err := godotenv.Load()
	if err != nil {
		 return err
	}

	Config.Token = os.Getenv("TOKEN")
	Config.BotPrefix = os.Getenv("BOT_PREFIX")
	Config.Undercover = os.Getenv("UNDERCOVER_WORDS")

	return nil
}
