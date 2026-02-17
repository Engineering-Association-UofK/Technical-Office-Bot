package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Engineering-Association-UofK/Technical-Office-Bot/internal/models"
	"github.com/Engineering-Association-UofK/Technical-Office-Bot/internal/service"
)

type FeedbackHandler struct {
	service service.FeedbackService
}

func NewFeedbackHandler(fbService *service.FeedbackService) *FeedbackHandler {
	return &FeedbackHandler{service: *fbService}
}

// HTTP handler
func (fbH *FeedbackHandler) HandleFeedbackRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.FeedbackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	feedback := models.FeedbackModel{
		Name:    req.Name,
		Email:   req.Email,
		Phone:   req.Phone,
		Message: req.Message,
	}

	id, err := fbH.service.Save(feedback)
	if err != nil {
		slog.Error("Failed to save feedback: " + err.Error())
		http.Error(w, "Failed to save feedback", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"id":     id,
	})
}
