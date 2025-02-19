package config

import (
	"os"

	"github.com/Safmica/undercover-bot/models"
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

	return nil
}
