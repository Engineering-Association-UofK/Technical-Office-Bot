package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/shirou/gopsutil/v4/net"
	"github.com/shirou/gopsutil/v4/process"
)

type Health struct {
	Groups []string `json:"groups"`
	Status string   `json:"status"`
}

func (sh *SystemHealth) CheckProcess() {
	p, err := findBackend(8080)
	if err != nil {
		slog.Error("Failed to fetch process info: "+err.Error(), "port", 8080)
		sh.AppProcess.IsRunning = false
		return
	}
	sh.AppProcess.IsRunning = true

	c, err := p.CPUPercent()
	if err != nil {
		slog.Error("Error reading process CPU usage: " + err.Error())
	} else {
		sh.AppProcess.CpuPercent = c
	}

	m, err := p.MemoryInfo()
	if err != nil {
		slog.Error("Error reading process CPU usage: " + err.Error())
	} else {
		sh.AppProcess.MemoryUsed = m.RSS
	}

}

func findBackend(targetPort uint32) (*process.Process, error) {
	conns, err := net.Connections("all")
	if err != nil {
		return nil, err
	}

	var pid int32

	for _, c := range conns {
		if c.Laddr.Port == targetPort {
			pid = c.Pid
		}
	}

	if pid == 0 {
		return nil, fmt.Errorf("Process with port not found")
	}

	procs, _ := process.Processes()
	for _, p := range procs {
		if pid == p.Pid {
			return p, nil
		}
	}

	return nil, fmt.Errorf("backend not found")
}

func (sh *SystemHealth) PingBackend() {
	resp, err := http.Get("http://localhost:8080/actuator/health")
	if err != nil {
		slog.Error("Error fetching actuator in Backend: "+err.Error(), "port", 8080, "endpoint", "/actuator/health")
		sh.IsResponsive = false
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Error Reading actuator request body: " + err.Error())
		sh.IsResponsive = false
		return
	}

	var health Health
	err = json.Unmarshal(body, &health)
	if err != nil {
		slog.Error("Error unmarshaling actuator request body: " + err.Error())
		sh.IsResponsive = false
		return
	} else {
		if health.Status == "UP" {
			sh.IsResponsive = true
		} else {
			sh.IsResponsive = false
		}
	}
}
