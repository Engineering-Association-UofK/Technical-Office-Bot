package repository

import (
	"fmt"

	"github.com/abdulrahim-m/Technical-Office-Bot/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramUserRepo struct {
	BaseRepo[models.TelegramUser]
}

func (tr *TelegramUserRepo) Save(tu *tgbotapi.User) (int64, error) {
	_, err := tr.Persist(`INSERT INTO `+tr.TableName+` (telegram_id, username, first_name, locale) 
		VALUES (?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			username = VALUES(username),
			first_name = VALUES(first_name),
			is_bot_blocked = FALSE;`,
		tu.ID, tu.UserName, tu.FirstName, tu.LanguageCode)
	if err != nil {
		return 0, err
	}

	return tu.ID, nil
}

func (tr *TelegramUserRepo) FindById(id int64) (models.TelegramUser, error) {
	var entry models.TelegramUser
	query := fmt.Sprintf("SELECT * FROM %s WHERE telegram_id = ?", tr.TableName)
	err := tr.DB.Get(&entry, query, id)
	return entry, err
}
