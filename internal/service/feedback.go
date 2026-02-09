package service

import (
	"fmt"
	"time"

	"github.com/abdulrahim-m/Technical-Office-Bot/internal/clients/telegram"
	"github.com/abdulrahim-m/Technical-Office-Bot/internal/models"
	"github.com/abdulrahim-m/Technical-Office-Bot/internal/repository"
)

type FeedbackService struct {
	Repo repository.FeedbackRepo
	Bot  *telegram.TelegramBot
}

func (fbS *FeedbackService) GetAll() ([]models.FeedbackModel, error) {
	return fbS.Repo.GetAll()
}

func (fbS *FeedbackService) FineById(id int) (models.FeedbackModel, error) {
	return fbS.Repo.FindById(id)
}

func (fbS *FeedbackService) Save(fb models.FeedbackModel) (int64, error) {
	fb.CreatedAt = time.Now()

	if fb.Message == "" {
		return 0, fmt.Errorf("invalid feedback data")
	}

	if fb.Name == "" {
		fb.Name = "Anonymous"
	}

	id, err := fbS.Repo.Save(fb)
	if err != nil {
		return 0, err
	}

	notificationMsg := fmt.Sprintf("New Feedback Received:\nName: %s\nEmail: %s\nPhone: %s\nMessage: %s",
		fb.Name, fb.Email, fb.Phone, fb.Message)
	fbS.Bot.NotifyAdmin(notificationMsg)

	return id, nil
}
