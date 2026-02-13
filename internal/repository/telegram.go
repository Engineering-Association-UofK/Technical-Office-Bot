package repository

import (
	"fmt"
	"log/slog"

	"github.com/abdulrahim-m/Technical-Office-Bot/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramRepo struct {
	Tuser        BaseRepo[models.TelegramUser]
	Tinteraction BaseRepo[models.TelegramInteraction]
}

func (tr *TelegramRepo) InteractionSave(tm *tgbotapi.Message) (int64, error) {
	_, err := tr.Tuser.Persist(`INSERT INTO `+tr.Tuser.TableName+` (telegram_id, username, first_name, locale, preferences) 
		VALUES (?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			username = VALUES(username),
			first_name = VALUES(first_name),
			locale = VALUES(locale),
			is_bot_blocked = FALSE;`,
		tm.From.ID, tm.From.UserName, tm.From.FirstName, tm.From.LanguageCode, models.NewPreferences())
	if err != nil {
		return 0, err
	}

	_, err = tr.Tuser.Persist(`INSERT INTO `+tr.Tinteraction.TableName+` (telegram_user_id, telegram_chat_id, message, created_at) 
		VALUES (?, ?, ?, ?)`,
		tm.From.ID, tm.Chat.ID, tm.Text, tm.Time())
	if err != nil {
		return 0, err
	}

	return tm.From.ID, nil
}

func (tr *TelegramRepo) UpdatePreferences(tu models.TelegramUser) error {
	if _, err := tr.FindById(tu.TelegramID); err != nil {
		slog.Error("Telegram user Was not found: "+err.Error(), "ID", tu.TelegramID, "Context", "Update Preferences")
	}

	query := fmt.Sprintf(`UPDATE %s SET preferences = ? WHERE telegram_id = ?`, tr.Tuser.TableName)
	_, err := tr.Tuser.DB.Exec(query, tu.Preferences, tu.TelegramID)
	return err
}

func (tr *TelegramRepo) FindById(id int64) (models.TelegramUser, error) {
	var entry models.TelegramUser
	query := fmt.Sprintf("SELECT * FROM %s WHERE telegram_id = ?", tr.Tuser.TableName)
	err := tr.Tuser.DB.Get(&entry, query, id)
	return entry, err
}

func (tr *TelegramRepo) FindNotifyEnabled() ([]models.TelegramUser, error) {
	var entries []models.TelegramUser
	query := fmt.Sprintf(`SELECT * FROM %s WHERE JSON_EXTRACT(preferences, '$.notify') = CAST(true AS JSON)`, tr.Tuser.TableName)
	err := tr.Tuser.DB.Select(&entries, query)
	return entries, err
}

func (tr *TelegramRepo) FindTechnicalAdmins() ([]models.TelegramUser, error) {
	var entries []models.TelegramUser
	query := fmt.Sprintf(`SELECT * FROM %s WHERE technical_admin = true`, tr.Tuser.TableName)
	err := tr.Tuser.DB.Select(&entries, query)
	return entries, err
}
