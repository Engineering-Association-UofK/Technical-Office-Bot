package main

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/Engineering-Association-UofK/Technical-Office-Bot/cmd/routes"
	"github.com/Engineering-Association-UofK/Technical-Office-Bot/internal/clients/telegram"
	"github.com/Engineering-Association-UofK/Technical-Office-Bot/internal/config"
	"github.com/Engineering-Association-UofK/Technical-Office-Bot/internal/database"
	"github.com/Engineering-Association-UofK/Technical-Office-Bot/internal/handler"
	"github.com/Engineering-Association-UofK/Technical-Office-Bot/internal/service"
)

var (
	fbHandler *handler.FeedbackHandler
	hHandler  *handler.HealthHandler
	hhHandler *handler.HelperHandler
)

func main() {
	Init()
	SetupHandlers()

	// Set up HTTP server and map endpoints
	routes.HttpStart(fbHandler, hHandler, hhHandler)
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
	db, err := database.NewMySQLConnection(fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", config.App.DBUser, config.App.DBPassword, config.App.DBHost, config.App.DBPort, config.App.DBName))
	if err != nil {
		panic("Error creating database connection: " + err.Error())
	}

	// Set up feedback
	fbService := service.NewFeedbackService(db, notificationChannel)
	fbHandler = handler.NewFeedbackHandler(fbService)

	// Start the telegram bot
	_, err = telegram.TelegramInit(config.App.TelegramToken, db, fbService, notificationChannel)
	if err != nil {
		slog.Error("Error starting telegram service: " + err.Error())
	}

	health, err := service.NewSystemHealth(sysHealthIntervalUpdateChannel, admin)
	if err != nil {
		panic("Error starting system monitoring: " + err.Error())
	}
	hHandler = handler.NewHealthHandler(health)

	pdfService := service.NewPDFService(1)
	hhHandler = handler.NewHelper(*pdfService)
}
