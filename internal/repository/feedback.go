package repository

import (
	"github.com/abdulrahim-m/Technical-Office-Bot/internal/models"
)

type FeedbackRepo struct {
	BaseRepo[models.FeedbackModel]
}

func (fbR *FeedbackRepo) GetAll() ([]models.FeedbackModel, error) {
	var fbs []models.FeedbackModel
	err := fbR.DB.Select(&fbs, "SELECT * FROM feedback")
	return fbs, err
}

func (fbR *FeedbackRepo) FineById(id int) (models.FeedbackModel, error) {
	fb := models.FeedbackModel{}

	err := fbR.DB.Get(&fb, "SELECT * FROM feedback WHERE id = ?", id)
	return fb, err
}

func (fbR *FeedbackRepo) Persist(fb models.FeedbackModel) (int64, error) {
	result, err := fbR.DB.Exec("INSERT INTO feedback (name, email, phone, message) VALUES (?, ?, ?, ?)", fb.Name, fb.Email, fb.Phone, fb.Message)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}
