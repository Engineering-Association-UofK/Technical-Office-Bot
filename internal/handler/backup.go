package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Engineering-Association-UofK/Technical-Office-Bot/internal/models"
	"github.com/Engineering-Association-UofK/Technical-Office-Bot/internal/service/backup"
)

type BackupHandler struct {
	service *backup.BackupService
}

func NewBackupHandler(backupService *backup.BackupService) *BackupHandler {
	return &BackupHandler{service: backupService}
}

func (b *BackupHandler) HandleTriggerRestore(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	go b.service.Restore()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Restore triggered successfully. Check logs for details."))
}

func (b *BackupHandler) HandleTriggerBackup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	go b.service.TriggerBackup()

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Backup triggered successfully"))
}

func (b *BackupHandler) HandleGetBackupDetails(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	details, err := b.service.GetBackupDetails()
	if err != nil {
		http.Error(w, "Failed to fetch backup details: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(details)
}

func (b *BackupHandler) HandleChangeBackupInterval(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Interval int64 `json:"interval"` // interval in nanoseconds
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := b.service.SetBackupInterval(time.Duration(req.Interval)); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Backup interval updated successfully"))
}

func (b *BackupHandler) HandleChangeBackupTime(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Hour   int `json:"hour"`
		Minute int `json:"minute"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := b.service.SetBackupTime(req.Hour, req.Minute); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Backup time updated successfully"))
}

func (b *BackupHandler) TriggerBackup(ctx context.Context, input *struct{}) (*models.RespWithBody, error) {
	go b.service.TriggerBackup()

	return &models.RespWithBody{
		Body: models.DefaultResponse{
			Status:    "success",
			Message:   "Backup triggered successfully",
			Timestamp: time.Now(),
		},
	}, nil
}

func (b *BackupHandler) TriggerRestore(ctx context.Context, input *struct{}) (*models.RespWithBody, error) {
	go b.service.Restore()

	return &models.RespWithBody{
		Body: models.DefaultResponse{
			Status:    "success",
			Message:   "Restore triggered successfully. Check logs for details.",
			Timestamp: time.Now(),
		},
	}, nil
}

func (b *BackupHandler) GetBackupDetails(ctx context.Context, input *struct{}) (*models.RespWithBody, error) {
	details, err := b.service.GetBackupDetails()
	if err != nil {
		return nil, err
	}

	return &models.RespWithBody{
		Body: details,
	}, nil
}

func (b *BackupHandler) ChangeBackupInterval(ctx context.Context, input *models.ChangeBackupIntervalRequest) (*models.RespWithBody, error) {
	if err := b.service.SetBackupInterval(input.Interval); err != nil {
		return nil, err
	}

	return &models.RespWithBody{
		Body: models.DefaultResponse{
			Status:    "success",
			Message:   "Backup interval updated successfully",
			Timestamp: time.Now(),
		},
	}, nil
}

func (b *BackupHandler) ChangeBackupTime(ctx context.Context, input *models.ChangeBackupTimeRequest) (*models.RespWithBody, error) {
	// Assuming the request contains a time.Time where we extract Hour and Minute
	if err := b.service.SetBackupTime(input.Time.Hour(), input.Time.Minute()); err != nil {
		return nil, err
	}

	return &models.RespWithBody{
		Body: models.DefaultResponse{
			Status:    "success",
			Message:   "Backup time updated successfully",
			Timestamp: time.Now(),
		},
	}, nil
}
