package repository

import (
	"github.com/Engineering-Association-UofK/Technical-Office-Bot/internal/models"
)

type FeedbackRepo struct {
	BaseRepo[models.FeedbackModel]
}

func (fbR *FeedbackRepo) Save(fb models.FeedbackModel) (int64, error) {
	return fbR.Persist(`INSERT INTO feedback (name, email, phone, message) VALUES (?, ?, ?, ?)`, fb.Name, fb.Email, fb.Phone, fb.Message)
}
