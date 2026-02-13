package models

import "time"

type HealthLevel string

const (
	Healthy  HealthLevel = "healthy"
	Warning  HealthLevel = "warning"
	Critical HealthLevel = "critical"
)

// Overview
type HealthOverviewResponse struct {
	HealthLevel   HealthLevel `json:"healthLevel"`
	UptimeSeconds uint64      `json:"uptimeSeconds"`
	LastUpdated   time.Time   `json:"lastUpdated"`
}

// System
type CPUResponse struct {
	LoadPercent float64 `json:"loadPercent"`
	Cores       int     `json:"cores"`
}

type MemoryResponse struct {
	UsedBytes       uint64  `json:"usedBytes"`
	TotalBytes      uint64  `json:"totalBytes"`
	UsedPercent     float64 `json:"usedPercent"`
	SwapUsedBytes   uint64  `json:"swapUsedBytes"`
	SwapTotalBytes  uint64  `json:"swapTotalBytes"`
	SwapUsedPercent float64 `json:"swapUsedPercent"`
}

type DiskResponse struct {
	UsedBytes   uint64  `json:"usedBytes"`
	TotalBytes  uint64  `json:"totalBytes"`
	UsedPercent float64 `json:"usedPercent"`
}

type SystemHealthResponse struct {
	CPU           CPUResponse    `json:"cpu"`
	Memory        MemoryResponse `json:"memory"`
	Disk          DiskResponse   `json:"disk"`
	UptimeSeconds uint64         `json:"uptimeSeconds"`
}

// App Process
type AppHealthResponse struct {
	IsRunning       bool    `json:"isRunning"`
	CPUPercent      float64 `json:"cpuPercent"`
	MemoryUsedBytes uint64  `json:"memoryUsedBytes"`
	MemoryUsedMB    uint64  `json:"memoryUsedMB"`
	Status          string  `json:"status"`
}

// Metrics
type MetricPoint struct {
	Timestamp time.Time `json:"timestamp"`
	Value     float64   `json:"value"`
}

type MetricsResponse struct {
	Range           string        `json:"range"`
	IntervalSeconds int           `json:"intervalSeconds"`
	CPU             []MetricPoint `json:"cpu"`
	Memory          []MetricPoint `json:"memory"`
}
