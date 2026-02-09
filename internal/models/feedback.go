package models

import (
	"database/sql"
	"time"
)

type FeedbackModel struct {
	ID         int           `db:"id"`
	Name       string        `db:"name"`
	Email      string        `db:"email"`
	Phone      string        `db:"phone"`
	TelegramID sql.NullInt64 `db:"telegram_id"`
	Message    string        `db:"message"`
	CreatedAt  time.Time     `db:"created_at"`
}
type FeedbackRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Message string `json:"message"`
}
