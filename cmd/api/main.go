package main

import (
	"log"

	"github.com/abdulrahim-m/Technical-Office-Bot/internal/config"
	"github.com/abdulrahim-m/Technical-Office-Bot/internal/telegram"
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

func main() {
	// Load the .env file
	godotenv.Load()

	cfg := config.Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Unable to parse config: %v", err)
	}

	// Start the telegram bot
	tBot := telegram.TelegramInit(cfg.TelegramToken, cfg.AdminTelegramID)

	log.Println(tBot)
	log.Println("Done!!!")

	for {
	}
}
