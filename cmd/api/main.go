package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/abdulrahim-m/Technical-Office-Bot/internal/clients/telegram"
	"github.com/abdulrahim-m/Technical-Office-Bot/internal/config"
	"github.com/abdulrahim-m/Technical-Office-Bot/internal/database"
	"github.com/abdulrahim-m/Technical-Office-Bot/internal/handler"
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

	// Set up database
	db := database.NewMySQLConnection(fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName))

	// Set up feedback
	fbHandler := handler.NewFeedbackHandler(db, tBot)

	// Set up HTTP server and map endpoints
	http.HandleFunc("/api/v1/feedback", fbHandler.HandleFeedbackRequest)

	log.Println("Starting server on :8080...")
	go func() {
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	log.Println("Bot is running...")

	for {
	}
}
