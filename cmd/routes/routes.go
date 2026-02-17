package routes

import (
	"net/http"

	"github.com/abdulrahim-m/Technical-Office-Bot/internal/handler"
)

func SetupRoutes(fb *handler.FeedbackHandler, h *handler.HealthHandler) http.Handler {
	mux := http.NewServeMux()

	feedBackHandler := handler.Basic(http.HandlerFunc(fb.HandleFeedbackRequest))
	mux.Handle("/api/v1/feedback", feedBackHandler)

	healthWithAuth := handler.Protected(http.HandlerFunc(h.HandleHealthRequests))
	mux.Handle("/api/v1/health/{path}/{value...}", healthWithAuth)

	return mux
}
