package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/abdulrahim-m/Technical-Office-Bot/internal/clients/telegram"
	"github.com/abdulrahim-m/Technical-Office-Bot/internal/models"
	"github.com/abdulrahim-m/Technical-Office-Bot/internal/repository"
	"github.com/abdulrahim-m/Technical-Office-Bot/internal/service"
	"github.com/jmoiron/sqlx"
)

type FeedbackHandler struct {
	service service.FeedbackService
}

func NewFeedbackHandler(db *sqlx.DB, bot *telegram.TelegramBot) *FeedbackHandler {
	return &FeedbackHandler{
		service: service.FeedbackService{
			Repo: repository.FeedbackRepo{
				BaseRepo: repository.BaseRepo[models.FeedbackModel]{
					DB:        db,
					TableName: "feedback",
				},
			},
			Bot: bot,
		},
	}
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
		log.Println("Failed to save feedback:", err)
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
