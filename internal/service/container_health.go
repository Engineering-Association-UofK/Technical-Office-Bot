package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/Engineering-Association-UofK/Technical-Office-Bot/internal/config"
	"github.com/Engineering-Association-UofK/Technical-Office-Bot/internal/models"
)

type ContainerHealthService struct {
	promURL    string
	httpClient *http.Client
	jobName    string

	// Cache state
	mu             sync.RWMutex
	appHealth      models.AppHealthResponse
	sysHealth      models.SystemHealthResponse
	metricsHistory models.MetricsHistoryResponse
}

func NewContainerHealthService() *ContainerHealthService {
	return &ContainerHealthService{
		promURL:    config.App.PrometheusHost,
		jobName:    "sea-api",
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}
}

func (s *ContainerHealthService) StartBackgroundUpdater(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	slog.Info("ContainerHealthService background updater started")

	go func() {
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				slog.Info("ContainerHealthService background updater stopped")
				return
			case <-ticker.C:
				s.RefreshData()
			}
		}
	}()

	s.RefreshData()
}

func (s *ContainerHealthService) RefreshData() {
	appUp := s.queryInstant(fmt.Sprintf(`up{job="%s"}`, s.jobName))
	appMem := s.queryInstant(fmt.Sprintf(`process_resident_memory_bytes{job="%s"}`, s.jobName))
	appCpu := s.queryInstant(fmt.Sprintf(`rate(process_cpu_seconds_total{job="%s"}[1m]) * 100`, s.jobName))

	s.mu.Lock()
	defer s.mu.Unlock()

	status := "Down"
	if appUp == 1 {
		status = "Running"
	}

	s.appHealth = models.AppHealthResponse{
		IsRunning:       appUp == 1,
		Status:          status,
		CPUPercent:      appCpu,
		MemoryUsedBytes: uint64(appMem),
		MemoryUsedMB:    uint64(appMem) / 1024 / 1024,
	}
}

func (s *ContainerHealthService) GetAppHealth() models.AppHealthResponse {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.appHealth
}

func (s *ContainerHealthService) queryInstant(query string) float64 {
	reqURL := fmt.Sprintf("%s/api/v1/query?query=%s", s.promURL, url.QueryEscape(query))

	resp, err := s.httpClient.Get(reqURL)
	if err != nil {
		slog.Error("Failed to reach Prometheus", "error", err)
		return 0
	}
	defer resp.Body.Close()

	var pResp models.PromResponse
	if err := json.NewDecoder(resp.Body).Decode(&pResp); err != nil || pResp.Status != "success" {
		slog.Error("Failed to decode Prometheus response", "error", err)
		return 0
	}

	// Safely extract the value from the array structure: [1714830000, "0.05"]
	if len(pResp.Data.Result) > 0 && len(pResp.Data.Result[0].Value) == 2 {
		valStr, ok := pResp.Data.Result[0].Value[1].(string)
		if ok {
			valFloat, _ := strconv.ParseFloat(valStr, 64)
			return valFloat
		}
	}

	return 0
}
