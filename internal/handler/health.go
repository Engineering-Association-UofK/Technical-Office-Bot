package handler

import (
	"context"
	"time"

	"github.com/Engineering-Association-UofK/Technical-Office-Bot/internal/models"
	"github.com/Engineering-Association-UofK/Technical-Office-Bot/internal/service"
	"github.com/Engineering-Association-UofK/Technical-Office-Bot/internal/utils"
)

type HealthHandler struct {
	System *service.SystemHealth
}

func NewHealthHandler(health *service.SystemHealth) *HealthHandler {
	return &HealthHandler{System: health}
}

func (hh *HealthHandler) GetOverview(ctx context.Context, input *struct{}) (*models.RespWithBody, error) {
	overview := hh.createOverview()
	return &models.RespWithBody{
		Body: overview,
	}, nil
}

func (hh *HealthHandler) GetSystemHealth(ctx context.Context, input *struct{}) (*models.RespWithBody, error) {
	system := hh.createSystemHealthDetails()
	return &models.RespWithBody{
		Body: system,
	}, nil
}

func (hh *HealthHandler) GetAppHealth(ctx context.Context, input *struct{}) (*models.RespWithBody, error) {
	app := hh.createAppHealthDetails()
	return &models.RespWithBody{
		Body: app,
	}, nil
}

func (hh *HealthHandler) createAppHealthDetails() models.AppHealthResponse {
	app := hh.System.AppProcess

	status := "DOWN"
	if app.IsRunning && hh.System.IsResponsive {
		status = "UP"
	}

	return models.AppHealthResponse{
		IsRunning:       app.IsRunning,
		CPUPercent:      app.CpuPercent,
		MemoryUsedBytes: app.MemoryUsed,
		MemoryUsedMB:    app.MemoryUsed / 1024 / 1024,
		Status:          status,
	}
}

func (hh *HealthHandler) createOverview() models.HealthOverviewResponse {
	var health models.HealthLevel
	if !hh.System.IsResponsive {
		health = models.Critical
	} else {
		if utils.MB_Conv(hh.System.Status.Memory.Max-hh.System.Status.Memory.Current) < 150 {
			health = models.Warning
		} else {
			health = models.Healthy
		}
	}

	return models.HealthOverviewResponse{
		HealthLevel:   health,
		UptimeSeconds: hh.System.Status.Uptime,
		LastUpdated:   time.Now(),
	}
}

func (hh *HealthHandler) createSystemHealthDetails() models.SystemHealthResponse {
	c := models.CPUResponse{
		LoadPercent: hh.System.Status.CPU.Load,
		Cores:       hh.System.Status.CPU.CoreNum,
	}

	used := hh.System.Status.Memory.Current
	max := hh.System.Status.Memory.Max
	swapUsed := hh.System.Status.Memory.CurrentSwap
	swapMax := hh.System.Status.Memory.MaxSwap
	percent := (float64(used) / float64(max))
	var swapPercent float64
	if swapMax != 0 {
		swapPercent = (float64(swapUsed) / float64(swapMax))
	}

	m := models.MemoryResponse{
		UsedBytes:       used,
		TotalBytes:      max,
		UsedPercent:     percent,
		SwapUsedBytes:   swapUsed,
		SwapTotalBytes:  swapMax,
		SwapUsedPercent: swapPercent,
	}

	used = hh.System.Status.Disk.Current
	max = hh.System.Status.Disk.Max
	percent = float64(used) / float64(max)

	d := models.DiskResponse{
		UsedBytes:   used,
		TotalBytes:  max,
		UsedPercent: percent,
	}

	return models.SystemHealthResponse{
		CPU:           c,
		Memory:        m,
		Disk:          d,
		UptimeSeconds: hh.System.Status.Uptime,
	}
}
