package routes

import (
	"log/slog"
	"net/http"

	"github.com/Engineering-Association-UofK/Technical-Office-Bot/internal/config"
	"github.com/Engineering-Association-UofK/Technical-Office-Bot/internal/handler"
	"github.com/rs/cors"
)

func HttpStart(fb *handler.FeedbackHandler, h *handler.HealthHandler) {
	mux := http.NewServeMux()

	feedBackHandler := handler.Basic(http.HandlerFunc(fb.HandleFeedbackRequest))
	mux.Handle("/api/v1/feedback", feedBackHandler)

	healthWithAuth := handler.Protected(http.HandlerFunc(h.HandleHealthRequests))
	mux.Handle("/api/v1/health/{path...}", healthWithAuth)

	slog.Info("Starting server on :" + config.App.Port)

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
