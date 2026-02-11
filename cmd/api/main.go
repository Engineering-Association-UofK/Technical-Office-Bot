package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/abdulrahim-m/Technical-Office-Bot/internal/clients/server"
	"github.com/abdulrahim-m/Technical-Office-Bot/internal/clients/telegram"
	"github.com/abdulrahim-m/Technical-Office-Bot/internal/config"
	"github.com/abdulrahim-m/Technical-Office-Bot/internal/database"
	"github.com/abdulrahim-m/Technical-Office-Bot/internal/handler"
	"github.com/abdulrahim-m/Technical-Office-Bot/internal/service"
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

	notificationChannel := make(chan string, 25)
	sysHelthIntervalUpdateChannel := make(chan time.Duration, 1)

	_ = server.NewSystemHealth(sysHelthIntervalUpdateChannel)

	// Set up database
	db := database.NewMySQLConnection(fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName))

	// Set up feedback
	fbService := service.NewFeedbackService(db, notificationChannel)
	fbHandler := handler.NewFeedbackHandler(fbService)

	// Start the telegram bot
	_ = telegram.TelegramInit(cfg.TelegramToken, cfg.AdminTelegramID, db, fbService, notificationChannel)

	// Set up HTTP server and map endpoints
	http.HandleFunc("/api/v1/feedback", fbHandler.HandleFeedbackRequest)

	log.Printf("Starting server on :%s...\n", cfg.Port)
	go func() {
		if err := http.ListenAndServe(":"+cfg.Port, nil); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	log.Println("Bot is running...")

	for {
	}
}
