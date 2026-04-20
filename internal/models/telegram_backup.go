package models

import "time"

type Telegrambackup struct {
	ID         int    `db:"id"`
	FileID     string `db:"telegram_file_id"`
	BackupDate string `db:"backup_date"`
	PartNum    int    `db:"part_number"`
}

type BackupDetailsResponse struct {
	LatestBackupDate string        `json:"latestBackupDate"`
	BackupTime       time.Time     `json:"backupTime"`
	BackupInterval   time.Duration `json:"backupInterval"`
}

type ChangeBackupIntervalRequest struct {
	Interval time.Duration `json:"interval"`
}

type RestoreWithDateRequest struct {
	Date time.Time `json:"date"`
}

type ChangeBackupTimeRequest struct {
	Time time.Time `json:"time"`
}
