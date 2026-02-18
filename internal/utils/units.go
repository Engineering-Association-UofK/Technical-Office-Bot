package utils

import (
	"fmt"
)

func KB(u uint64) string {
	return fmt.Sprintf("%.2f", float64(u)/1024)
}

func KB_Conv(u uint64) float64 {
	return float64(u) / 1024
}

func MB(u uint64) string {
	return fmt.Sprintf("%.2f", float64(u)/1_048_576)
}

func MB_Conv(u uint64) float64 {
	return float64(u) / 1_048_576
}

func GB(u uint64) string {
	return fmt.Sprintf("%.2f", float64(u)/1_073_741_824)
}

func GB_Conv(u uint64) float64 {
	return float64(u) / 1_073_741_824
}

func Digit2(f float64) string {
	return fmt.Sprintf("%.2f", f)
}

func Digit3(f float64) string {
	return fmt.Sprintf("%.3f", f)
}
