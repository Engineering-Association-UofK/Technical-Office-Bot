package models

import (
	"time"
)

type TelegramUser struct {
	TelegramID      int64     `db:"telegram_id"`
	Username        string    `db:"username"`
	FirstName       string    `db:"first_name"`
	Locale          string    `db:"locale"`
	IsBotBlocked    bool      `db:"is_bot_blocked"`
	NotifyPromotion bool      `db:"notify_promotion"`
	CreatedAt       time.Time `db:"created_at"`
}
