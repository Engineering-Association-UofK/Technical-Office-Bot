package models

import (
	"time"
)

// --- Frontend Response Models ---

type AppHealthResponse struct {
	IsRunning       bool    `json:"isRunning"`
	Status          string  `json:"status"`
	CPUPercent      float64 `json:"cpuPercent"`
	MemoryUsedBytes uint64  `json:"memoryUsedBytes"`
	MemoryUsedMB    uint64  `json:"memoryUsedMB"`
}

type SystemHealthResponse struct {
	CPUPercent       float64 `json:"cpuPercent"`
	MemoryUsedBytes  uint64  `json:"memoryUsedBytes"`
	MemoryTotalBytes uint64  `json:"memoryTotalBytes"`
	DiskUsedBytes    uint64  `json:"diskUsedBytes"`
	DiskTotalBytes   uint64  `json:"diskTotalBytes"`
	UptimeSeconds    uint64  `json:"uptimeSeconds"`
}

type MetricPoint struct {
	Timestamp time.Time `json:"timestamp"`
	Value     float64   `json:"value"`
}

type MetricsHistoryResponse struct {
	AppCPU    []MetricPoint `json:"appCpu"`
	AppMemory []MetricPoint `json:"appMemory"`
}

// --- Prometheus JSON Parsing Models ---
// Prometheus returns values as an array of mixed types: [1714830000.123, "0.05"] (Timestamp, StringValue)

type PromResponse struct {
	Status string   `json:"status"`
	Data   promData `json:"data"`
}

type promData struct {
	ResultType string       `json:"resultType"`
	Result     []promResult `json:"result"`
}

type promResult struct {
	Metric map[string]string `json:"metric"`
	Value  []interface{}     `json:"value"`  // For instant queries
	Values [][]interface{}   `json:"values"` // For range queries
}
