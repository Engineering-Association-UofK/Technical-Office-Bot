package models

import "database/sql"

type Admin struct {
	ID         int           `db:"id"`
	Name       string        `db:"name"`
	Email      string        `db:"email"`
	TelegramID sql.NullInt64 `db:"telegram_id"`
	DiscordID  string        `db:"discord_id"`
}
