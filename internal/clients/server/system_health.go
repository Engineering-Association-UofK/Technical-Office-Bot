package server

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
)

type SystemHealth struct {
	Status     Status
	AppProcess AppProcess
	Interval   time.Duration

	IsResponsive bool

	IntervalChannel <-chan time.Duration
}

func NewSystemHealth(channel <-chan time.Duration) (*SystemHealth, error) {
	sh := SystemHealth{
		Status: Status{
			Memory: Memory{},
			Disk:   Disk{},
			CPU:    CPU{},
		},
		Interval:        time.Second * 5,
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
		formatUptime(sh.Status.Uptime),
		sh.Status.CPU.CoreNum,
		MB(sh.Status.Memory.Current),
		MB(sh.Status.Memory.Max),
		MB(sh.Status.Memory.CurrentSwap),
		MB(sh.Status.Memory.MaxSwap),
		Digit3(sh.Status.CPU.Load),
		GB(sh.Status.Disk.Current),
		GB(sh.Status.Disk.Max),
		sh.AppProcess.IsRunning,
		Digit3(sh.AppProcess.CpuPercent),
		MB(sh.AppProcess.MemoryUsed),
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
