package server

import (
	"fmt"
	"strings"
)

func formatUptime(seconds uint64) string {
	days := seconds / 86400
	hours := (seconds % 86400) / 3600
	minutes := (seconds % 3600) / 60
	secs := seconds % 60

	parts := []string{}

	if days > 0 {
		parts = append(parts, fmt.Sprintf("%dd", days))
	}
	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%dh", hours))
	}
	if minutes > 0 {
		parts = append(parts, fmt.Sprintf("%dm", minutes))
	}
	if secs > 0 || len(parts) == 0 {
		parts = append(parts, fmt.Sprintf("%ds", secs))
	}

	return strings.Join(parts, " ")
}

func kb(u uint64) string {
	return fmt.Sprintf("%.2f", float64(u)/1000)
}

func mb(u uint64) string {
	return fmt.Sprintf("%.2f", float64(u)/1000000)
}

func gb(u uint64) string {
	return fmt.Sprintf("%.2f", float64(u)/1000000000)
}
