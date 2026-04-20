package routes

import (
	"log/slog"
	"net/http"

	"github.com/Engineering-Association-UofK/Technical-Office-Bot/internal/config"
	"github.com/Engineering-Association-UofK/Technical-Office-Bot/internal/handler"
	"github.com/rs/cors"
)

var (
	Fbh *handler.FeedbackHandler
	Hh  *handler.HealthHandler
	Bh  *handler.BackupHandler
)

func HttpStart() {
	mux := http.NewServeMux()

	////////////////////
	////  Feedback  ////
	////////////////////

	feedBackHandler := handler.Basic(http.HandlerFunc(Fbh.HandleFeedbackRequest))
	mux.Handle("/api/v1/feedback", feedBackHandler)

	////////////////////
	////   Health   ////
	////////////////////

	healthWithAuth := handler.Protected(http.HandlerFunc(Hh.HandleHealthRequests))
	mux.Handle("/api/v1/health/{path...}", healthWithAuth)

	////////////////////
	////   Backup   ////
	////////////////////

	restore := handler.Protected(http.HandlerFunc(Bh.HandleTriggerRestore))
	mux.Handle("/api/v1/restore", restore)

	backup := handler.Protected(http.HandlerFunc(Bh.HandleTriggerBackup))
	mux.Handle("/api/v1/backup", backup)

	backupDetails := handler.Protected(http.HandlerFunc(Bh.HandleGetBackupDetails))
	mux.Handle("/api/v1/backup/details", backupDetails)

	backupInterval := handler.Protected(http.HandlerFunc(Bh.HandleChangeBackupInterval))
	mux.Handle("/api/v1/backup/interval", backupInterval)

	backupTime := handler.Protected(http.HandlerFunc(Bh.HandleChangeBackupTime))
	mux.Handle("/api/v1/backup/time", backupTime)

	////////////////////
	////   Config   ////
	////////////////////

	slog.Info("Routes registered successfully")
	slog.Info("Starting server on :" + config.App.Port)

	// CORS
	corsOptions := cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"X-Requested-With", "Content-Type", "Authorization"},
	}

	handler := cors.New(corsOptions).Handler(mux)

	if err := http.ListenAndServe(":"+config.App.Port, handler); err != nil {
		panic("Server failed: " + err.Error())
	}
}
