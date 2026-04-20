package repository

import (
	"github.com/Engineering-Association-UofK/Technical-Office-Bot/internal/models"
	"github.com/jmoiron/sqlx"
)

type TelegramBackupRepo struct {
	BaseRepo[models.Telegrambackup]
}

func (tbr *TelegramBackupRepo) Save(backup models.Telegrambackup) (int64, error) {
	return tbr.Persist(`INSERT INTO telegram_backups (telegram_file_id, backup_date, part_number) VALUES (?, ?, ?)`,
		backup.FileID, backup.BackupDate, backup.PartNum)
}

// Save with transaction
func (tbr *TelegramBackupRepo) SaveTx(tx *sqlx.Tx, backup models.Telegrambackup) (int64, error) {
	return tbr.PersistTx(tx, `INSERT INTO telegram_backups (telegram_file_id, backup_date, part_number) VALUES (?, ?, ?)`,
		backup.FileID, backup.BackupDate, backup.PartNum)
}

func (tbr *TelegramBackupRepo) FindByDate(date string) ([]models.Telegrambackup, error) {
	var entries []models.Telegrambackup
	query := `SELECT * FROM telegram_backups WHERE backup_date = ? ORDER BY part_number ASC`
	err := tbr.DB.Select(&entries, query, date)
	return entries, err
}

func (tbr *TelegramBackupRepo) FindLatestDate() (string, error) {
	var date string
	query := `SELECT backup_date FROM telegram_backups ORDER BY backup_date DESC LIMIT 1`
	return date, tbr.DB.Get(&date, query)
}
