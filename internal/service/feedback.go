package service

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/abdulrahim-m/Technical-Office-Bot/internal/models"
	"github.com/abdulrahim-m/Technical-Office-Bot/internal/repository"
	"github.com/jmoiron/sqlx"
)

type FeedbackService struct {
	Repo   repository.FeedbackRepo
	notify chan<- string
}

func NewFeedbackService(db *sqlx.DB, notificationChannel chan<- string) *FeedbackService {
	return &FeedbackService{
		Repo: repository.FeedbackRepo{
			BaseRepo: repository.BaseRepo[models.FeedbackModel]{
				DB:        db,
				TableName: "feedback",
			},
		},
		notify: notificationChannel,
	}
}

func (fbS *FeedbackService) GetAll() ([]models.FeedbackModel, error) {
	return fbS.Repo.GetAll()
}

func (fbS *FeedbackService) FineById(id int64) (models.FeedbackModel, error) {
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
	fbS.notify <- notificationMsg

	return id, nil
}

func (fbS *FeedbackService) NotifyFeedback(name string, ID int64, message string) (int64, error) {
	fb := models.FeedbackModel{
		Name:       name,
		TelegramID: sql.NullInt64{Int64: ID, Valid: ID != 0},
		Message:    message,
	}

	id, err := fbS.Repo.Save(fb)
	if err != nil {
		return 0, err
	}

	notificationMsg := fmt.Sprintf("New Telegram Feedback Received:\n\nName: %s\nTelegram ID: %v\nMessage: %s",
		fb.Name, fb.TelegramID.Int64, fb.Message)
	fbS.notify <- notificationMsg

	return id, nil
}
