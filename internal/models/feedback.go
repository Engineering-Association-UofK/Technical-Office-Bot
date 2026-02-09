package models

import "time"

type FeedbackModel struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	Phone     string    `db:"phone"`
	Message   string    `db:"message"`
	CreatedAt time.Time `db:"created_at"`
}

type FeedbackRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Message string `json:"message"`
}
