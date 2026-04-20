package backup

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/Engineering-Association-UofK/Technical-Office-Bot/internal/clients/telegram"
	"github.com/Engineering-Association-UofK/Technical-Office-Bot/internal/config"
	"github.com/Engineering-Association-UofK/Technical-Office-Bot/internal/models"
	"github.com/Engineering-Association-UofK/Technical-Office-Bot/internal/repository"
	"github.com/jmoiron/sqlx"
)

type BackupService struct {
	bot  *telegram.TelegramBot
	repo repository.TelegramBackupRepo

	stop chan struct{}

	// Backup

	isBackupRunning bool

	sourceDir string
	secretKey string

	BackupTime     time.Time
	BackupInterval time.Duration
}

func NewBackupService(db sqlx.DB, bot *telegram.TelegramBot) *BackupService {
	service := BackupService{
		bot: bot,
		repo: repository.TelegramBackupRepo{
			BaseRepo: repository.BaseRepo[models.Telegrambackup]{
				DB:        &db,
				TableName: "telegram_backups",
			},
		},
		stop:      make(chan struct{}),
		sourceDir: config.App.BackupDir,
		secretKey: config.App.BackupSecret,
		BackupTime: time.Date(
			time.Now().Year(),
			time.Now().Month(),
			time.Now().Day(),
			2,
			0,
			0,
			0,
			time.Local,
		),
		BackupInterval: time.Hour * 24,
	}

	go service.Start()

	return &service
}

func (b *BackupService) Stop() {
	close(b.stop)
}

func (b *BackupService) Restart() {
	if b.stop != nil {
		b.Stop()
	}
	b.stop = make(chan struct{})
	go b.Start()
}

func (b *BackupService) IsRunning() bool {
	return b.stop != nil
}

func (b *BackupService) SetBackupInterval(interval time.Duration) error {
	if interval <= 0 {
		return fmt.Errorf("invalid backup interval")
	} else if interval < time.Hour*3 {
		return fmt.Errorf("backup interval must be at least 3 hours")
	}
	if b.isBackupRunning {
		return fmt.Errorf("backup is running, cannot change interval during backup")
	}
	b.BackupInterval = interval
	b.Restart()
	return nil
}

func (b *BackupService) SetBackupTime(h int, m int) error {
	if h < 0 || h > 23 || m < 0 || m > 59 {
		return fmt.Errorf("invalid time")
	}
	if b.isBackupRunning {
		return fmt.Errorf("backup is running, cannot change backup time during backup")
	}
	b.BackupTime = time.Date(
		time.Now().Year(),
		time.Now().Month(),
		time.Now().Day(),
		h,
		m,
		0,
		0,
		time.Local,
	)
	b.Restart()
	return nil
}

func (b *BackupService) TriggerBackup() {
	b.Backup()
}

func (b *BackupService) GetBackupDetails() (models.BackupDetailsResponse, error) {
	date, err := b.repo.FindLatestDate()
	if err != nil {
		return models.BackupDetailsResponse{}, err
	}
	return models.BackupDetailsResponse{
		LatestBackupDate: date,
		BackupTime:       b.BackupTime,
		BackupInterval:   b.BackupInterval,
	}, nil
}

func (b *BackupService) Start() {
	var wait time.Duration
	if time.Now().After(b.BackupTime) {
		wait = time.Since(b.BackupTime)
	} else {
		wait = time.Until(b.BackupTime)
	}

	isDoneFirst := false
	var timer time.Ticker
	for !isDoneFirst {
		t := time.NewTicker(wait)

		<-t.C

		isDoneFirst = true
		timer = *time.NewTicker(b.BackupInterval)
		t.Stop()
		b.Backup()
	}

	for {
		select {
		case <-timer.C:
			b.Backup()
		case <-b.stop:
			timer.Stop()
			return
		}
	}
}

func (b *BackupService) Backup() {
	slog.Info("Starting backup...")
	b.isBackupRunning = true
	defer func() {
		b.isBackupRunning = false
	}()

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		fmt.Println("Backup failed:", err)
		return
	}

	pr, pw := io.Pipe()
	go CompressAndEncrypt(pw, b.sourceDir, b.secretKey, iv)

	chunkSize := int64(15 * 1024 * 1024)
	slog.Info("Fixed chuck size", "size", chunkSize)
	t := time.Now()
	partNum := 1

	tx, err := b.repo.DB.BeginTxx(context.Background(), nil)
	if err != nil {
		fmt.Println("Failed to start transaction:", err)
		return
	}
	defer tx.Rollback()

	for {
		name := fmt.Sprintf("backup_%d-%d-%d_part-%d.tar.gz.enc",
			t.Year(), t.Month(), t.Day(), partNum)

		outFile, err := os.Create("./" + name)
		if err != nil {
			fmt.Println("Failed to create chunk:", err)
			return
		}

		slog.Info("Backup chunk created", "part", partNum, "name", name)

		// IV is written only on the first part
		if partNum == 1 {
			outFile.Write(iv)
		}

		// Pull exactly 15MB out of the pipe (Telegram bot upload limit is 50MB, and download is 20MB, for automation we will go with 15MB)
		n, copyErr := io.CopyN(outFile, pr, chunkSize)
		outFile.Close()

		// If bytes were written, the chunk still needs to be processed
		if n > 0 || (partNum == 1 && n == 0) {
			dateStr := fmt.Sprintf("%d_%02d_%02d", t.Year(), t.Month(), t.Day())

			// Upload to Telegram
			fileID, err := b.bot.UploadChunk("./"+name, dateStr, partNum)
			if err != nil {
				fmt.Println("Telegram upload failed:", err)
				return
			}

			// Save record using the transaction
			_, err = b.repo.SaveTx(tx, models.Telegrambackup{
				FileID:     fileID,
				BackupDate: dateStr,
				PartNum:    partNum,
			})
			if err != nil {
				fmt.Println("DB insert failed:", err)
				return
			}
			os.Remove("./" + name)
			partNum++
		} else {
			os.Remove("./" + name)
		}

		if copyErr != nil {
			if copyErr == io.EOF {
				break
			}
			fmt.Println("Error writing chunk:", copyErr)
			return
		}
	}

	tx.Commit()
	fmt.Println("Backup successful!")
}

func (b *BackupService) Restore() error {
	slog.Info("Starting restore...")
	date, err := b.repo.FindLatestDate()
	if err != nil {
		return err
	}
	return b.RestoreWithDate(date)
}

func (b *BackupService) RestoreWithDate(date string) error {
	entries, err := b.repo.FindByDate(date)
	if err != nil || len(entries) == 0 {
		return fmt.Errorf("no backups found for latest date")
	}

	slog.Info("Opening chunks...")
	var readers []io.Reader
	for _, entry := range entries {
		stream, err := b.bot.GetFileStream(entry.FileID)
		if err != nil {
			return err
		}
		defer stream.Close()
		readers = append(readers, stream)
	}

	slog.Debug("Combining chunks...")
	multiReader := io.MultiReader(readers...)

	// Read IV from the beginning of the stream
	iv := make([]byte, aes.BlockSize)
	slog.Debug("Reading IV...")
	if _, err := io.ReadFull(multiReader, iv); err != nil {
		return fmt.Errorf("failed to read IV: %w", err)
	}

	slog.Debug("Starting Decryption...")
	block, err := aes.NewCipher([]byte(b.secretKey))
	if err != nil {
		return err
	}
	stream := cipher.NewCTR(block, iv)
	decReader := &cipher.StreamReader{S: stream, R: multiReader}

	slog.Debug("Starting Gzip...")
	gzReader, err := gzip.NewReader(decReader)
	if err != nil {
		return err
	}
	defer gzReader.Close()

	slog.Debug("Starting Tar...")
	tarReader := tar.NewReader(gzReader)

	// Delete everything in the directory
	slog.Debug("Cleaning target directory...")
	if err := os.RemoveAll(b.sourceDir); err != nil {
		return fmt.Errorf("failed to clean backup directory: %w", err)
	}
	if err := os.MkdirAll(b.sourceDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to recreate backup directory: %w", err)
	}

	slog.Debug("Unpacking...")
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break // Reached the end of the backup
		}
		if err != nil {
			return err
		}

		target := filepath.Join(b.sourceDir, header.Name)

		// Handle directories
		if header.FileInfo().IsDir() {
			os.MkdirAll(target, os.ModePerm)
			continue
		}

		// Ensure parent directories exist for files
		os.MkdirAll(filepath.Dir(target), os.ModePerm)

		// Write the actual file
		outFile, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
		if err != nil {
			return err
		}

		if _, err := io.Copy(outFile, tarReader); err != nil {
			outFile.Close()
			return err
		}
		outFile.Close()
	}

	slog.Info("Backup restored successfully!")
	return nil
}

// Helpers

func CompressAndEncrypt(pw *io.PipeWriter, sourceDir, secretKey string, iv []byte) {
	defer pw.Close()

	block, er := aes.NewCipher([]byte(secretKey))
	if er != nil {
		slog.Error(er.Error())
		return
	}
	stream := cipher.NewCTR(block, iv)

	encryptedWriter := &cipher.StreamWriter{S: stream, W: pw}
	defer encryptedWriter.Close()

	gzipWriter := gzip.NewWriter(encryptedWriter)
	defer gzipWriter.Close()

	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	err := filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := tar.FileInfoHeader(info, info.Name())
		if err != nil {
			return err
		}

		relPath, _ := filepath.Rel(filepath.Dir(sourceDir), path)
		header.Name = filepath.ToSlash(relPath)

		if err := tarWriter.WriteHeader(header); err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(tarWriter, file)
		return err
	})

	if err != nil {
		pw.CloseWithError(err)
	}
}
