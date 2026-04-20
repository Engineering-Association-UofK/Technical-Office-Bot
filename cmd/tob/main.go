package main

import (
	"log/slog"
	"time"

	"github.com/Engineering-Association-UofK/Technical-Office-Bot/cmd/routes"
	"github.com/Engineering-Association-UofK/Technical-Office-Bot/internal/clients/telegram"
	"github.com/Engineering-Association-UofK/Technical-Office-Bot/internal/config"
	"github.com/Engineering-Association-UofK/Technical-Office-Bot/internal/database"
	"github.com/Engineering-Association-UofK/Technical-Office-Bot/internal/handler"
	"github.com/Engineering-Association-UofK/Technical-Office-Bot/internal/service"
	"github.com/Engineering-Association-UofK/Technical-Office-Bot/internal/service/backup"
)

func main() {
	Init()

	slog.Info("Starting Technical Office Bot... 🤖")

	SetupHandlers()

	// Set up HTTP server and map endpoints
	routes.HttpStart()
}

func Init() {
	// Load Environment Variables
	if err := config.Load(); err != nil {
		slog.Error("Unable to parse config: " + err.Error())
		return
	}

	// Setup logging
	Log := config.NewMultiHandlerLog()
	slog.SetDefault(Log)
}

func SetupHandlers() {
	admin := service.NewAdminAccount()

	notificationChannel := make(chan string, 25)
	sysHealthIntervalUpdateChannel := make(chan time.Duration, 1)

	// Set up database
	db, err := database.NewMySQLConnection()
	if err != nil {
		panic("Error creating database connection: " + err.Error())
	}

	// Set up feedback
	fbService := service.NewFeedbackService(db, notificationChannel)
	routes.Fbh = handler.NewFeedbackHandler(fbService)

	health, err := service.NewSystemHealth(sysHealthIntervalUpdateChannel, admin)
	if err != nil {
		panic("Error starting system monitoring: " + err.Error())
	}
	routes.Hh = handler.NewHealthHandler(health)

	// Start the telegram bot
	bot, err := telegram.TelegramInit(config.App.TelegramToken, db, fbService, notificationChannel)
	if err != nil {
		slog.Error("Error starting telegram service: " + err.Error())
	}

	backupService := backup.NewBackupService(*db, bot)
	routes.Bh = handler.NewBackupHandler(backupService)
}
