package service

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/Engineering-Association-UofK/Technical-Office-Bot/internal/models"
	"github.com/Engineering-Association-UofK/Technical-Office-Bot/internal/utils"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/net"
	"github.com/shirou/gopsutil/v4/process"
)

type SystemHealth struct {
	Status     models.Status
	AppProcess models.AppProcess
	Interval   time.Duration

	AdminAccount *AdminAccount

	IsResponsive bool

	IntervalChannel <-chan time.Duration
}

func NewSystemHealth(channel <-chan time.Duration, admin *AdminAccount) (*SystemHealth, error) {
	sh := SystemHealth{
		Status: models.Status{
			Memory: models.Memory{},
			Disk:   models.Disk{},
			CPU:    models.CPU{},
		},
		Interval:        time.Second * 10,
		AdminAccount:    admin,
		IsResponsive:    true,
		IntervalChannel: channel,
	}

	coreNum, err := cpu.Counts(true)
	if err != nil {
		return nil, err
	}
	sh.Status.CPU.CoreNum = coreNum

	go sh.Listen()

	return &sh, nil
}

func (sh *SystemHealth) StatusTextSummery() string {
	return fmt.Sprintf(
		`
System Responsive: %v
		
System Uptime: %s.
Number of cores: %d.
Memory: %v/%vMB.
Swap: %v/%vMB.
CPU: %v%%. 
Disk: %v/%vGB.

Process Running: %v.
Process usage CPU: %v%%.
Process usage RAM: %v.
	`,
		sh.IsResponsive,
		utils.FormatSeconds(sh.Status.Uptime),
		sh.Status.CPU.CoreNum,
		utils.MB(sh.Status.Memory.Current),
		utils.MB(sh.Status.Memory.Max),
		utils.MB(sh.Status.Memory.CurrentSwap),
		utils.MB(sh.Status.Memory.MaxSwap),
		utils.Digit3(sh.Status.CPU.Load),
		utils.GB(sh.Status.Disk.Current),
		utils.GB(sh.Status.Disk.Max),
		sh.AppProcess.IsRunning,
		utils.Digit3(sh.AppProcess.CpuPercent),
		utils.MB(sh.AppProcess.MemoryUsed),
	)
}

func (sh *SystemHealth) Listen() {
	ticker := time.NewTicker(sh.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			sh.UpdateStatus()
			sh.CheckProcess()
			sh.PingBackend()
		case newInterval := <-sh.IntervalChannel:
			ticker.Reset(newInterval)
			sh.UpdateStatus()
		}
	}
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
	health, e, err := sh.AdminAccount.CheckHealth()
	if err != nil {
		slog.Error("Error fetching actuator in Backend: "+err.Error(), "port", 8080, "endpoint", "/actuator/health")
		sh.IsResponsive = false
		return
	}
	if e != nil {
		slog.Error("Pinging backend returned an exception", "Status", e.Status, "Message", e.Message, "Time", e.TimeStamp)
		sh.IsResponsive = false
		return
	}
	if health == nil {
		slog.Error("Pinging backend returned no data", "port", 8080, "endpoint", "/actuator/health")
		sh.IsResponsive = false
		return
	}

	if health.Status == "UP" {
		sh.IsResponsive = true
	} else {
		sh.IsResponsive = false
	}
}

func (sh *SystemHealth) UpdateStatus() {
	// RAM Stats
	v, err := mem.VirtualMemory()
	if err != nil {
		slog.Error("Error reading Virtual Memory stats: " + err.Error())
	} else {
		sh.Status.Memory.Update(
			v.Used,
			v.Total,
			v.SwapCached,
			v.SwapTotal,
		)
	}

	// CPU Stats
	c, err := cpu.Percent(0, false)
	if err != nil {
		slog.Error("Error reading CPU stats: " + err.Error())
	} else {
		sh.Status.CPU.Update(c[0])
	}

	// Disk Stats
	d, err := disk.Usage("/")
	if err != nil {
		slog.Error("Error reading Disk stats: " + err.Error())
	} else {
		sh.Status.Disk.Update(d.Used, d.Total)
	}

	// Uptime
	u, err := host.Uptime()
	if err != nil {
		slog.Error("Error reading Uptime: " + err.Error())
	} else {
		sh.Status.Uptime = u
	}

}
