package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type Preferences struct {
	Notify bool `json:"notify"`
}

func NewPreferences() Preferences {
	return Preferences{
		Notify: false,
	}
}

func (p Preferences) Value() (driver.Value, error) {
	return json.Marshal(p)
}

func (p *Preferences) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		// Handle cases where the driver returns a string instead of []byte
		if s, ok := value.(string); ok {
			b = []byte(s)
		} else {
			return fmt.Errorf("type assertion to []byte failed")
		}
	}
	return json.Unmarshal(b, p)
}

type TelegramUser struct {
	TelegramID     int64       `db:"telegram_id"`
	Username       string      `db:"username"`
	FirstName      string      `db:"first_name"`
	Locale         string      `db:"locale"`
	IsBotBlocked   bool        `db:"is_bot_blocked"`
	TechnicalAdmin bool        `db:"technical_admin"`
	Preferences    Preferences `db:"preferences"`
	CreatedAt      time.Time   `db:"created_at"`
}

type TelegramInteraction struct {
	ID             int       `db:"id"`
	TelegramUserID int64     `db:"telegram_user_id"`
	TelegramChatID int64     `db:"telegram_chat_id"`
	Message        string    `db:"message"`
	Type           *string   `db:"type"`
	CreatedAt      time.Time `db:"created_at"`
}
