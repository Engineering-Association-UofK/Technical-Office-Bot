package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/abdulrahim-m/Technical-Office-Bot/cmd/routes"
	"github.com/abdulrahim-m/Technical-Office-Bot/internal/clients/server"
	"github.com/abdulrahim-m/Technical-Office-Bot/internal/clients/telegram"
	"github.com/abdulrahim-m/Technical-Office-Bot/internal/config"
	"github.com/abdulrahim-m/Technical-Office-Bot/internal/database"
	"github.com/abdulrahim-m/Technical-Office-Bot/internal/handler"
	"github.com/abdulrahim-m/Technical-Office-Bot/internal/service"
)

func main() {
	// Load Environment Variables
	if err := config.Load(); err != nil {
		slog.Error("Unable to parse config: " + err.Error())
		return
	}

	// Setup logging
	Log := config.NewMultiHandlerLog()
	slog.SetDefault(Log)

	notificationChannel := make(chan string, 25)
	sysHealthIntervalUpdateChannel := make(chan time.Duration, 1)

	// Set up database
	db, err := database.NewMySQLConnection(fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", config.App.DBUser, config.App.DBPassword, config.App.DBHost, config.App.DBPort, config.App.DBName))
	if err != nil {
		slog.Error("Error creating database connection: " + err.Error())
		return
	}

	// Set up feedback
	fbService := service.NewFeedbackService(db, notificationChannel)
	fbHandler := handler.NewFeedbackHandler(fbService)

	// Start the telegram bot
	_, err = telegram.TelegramInit(config.App.TelegramToken, db, fbService, notificationChannel)
	if err != nil {
		slog.Error("Error starting telegram service: " + err.Error())
	}

	health, err := server.NewSystemHealth(sysHealthIntervalUpdateChannel)
	if err != nil {
		slog.Error("Error starting system monitoring: " + err.Error())
		return
	}
	hHandler := handler.NewHealthHandler(health)

	// Set up HTTP server and map endpoints
	router := routes.SetupRoutes(fbHandler, hHandler)

	slog.Info("Starting server on :" + config.App.Port)
	if err := http.ListenAndServe(":"+config.App.Port, router); err != nil {
		Log.Error("Server failed: " + err.Error())
		return
	}
}

// TODO: Use the Actuator structure
