package server

import (
	"fmt"
	"log"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/net"
	"github.com/shirou/gopsutil/v4/process"
)

type SystemHealth struct {
	Status     Status
	AppProcess AppProcess
	Interval   time.Duration

	IsResponcive bool

	IntervalChannel <-chan time.Duration
}

func NewSystemHealth(channel <-chan time.Duration) *SystemHealth {
	sh := SystemHealth{
		Status: Status{
			Memory: Memory{},
			Disk:   Disk{},
			CPU:    CPU{},
		},
		Interval:        time.Minute,
		IsResponcive:    true,
		IntervalChannel: channel,
	}

	coreNum, err := cpu.Counts(true)
	if err != nil {
		log.Panicln("Failed to get core count: ", err)
	}
	sh.Status.CPU.CoreNum = coreNum

	go sh.Listen()

	return &sh
}

func (sh *SystemHealth) StatusTextSummery() string {
	return fmt.Sprintf(
		`System Uptime: %s.
Number of cores: %d.
Mamory: %v/%vMB.
Swap: %v/%vMB.
CPU: %v%%. 
Disk: %v/%vGB.

Process Running: %v.
Process usage CPU: %v%%.
Process usage RAM: %v.
	`,
		formatUptime(sh.Status.Uptime),
		sh.Status.CPU.CoreNum,
		mb(sh.Status.Memory.Current),
		mb(sh.Status.Memory.Max),
		mb(sh.Status.Memory.CurremtSwap),
		mb(sh.Status.Memory.MaxSwap),
		sh.Status.CPU.Load,
		gb(sh.Status.Disk.Current),
		gb(sh.Status.Disk.Max),
		sh.AppProcess.IsRunning,
		sh.AppProcess.CpuPercent,
		sh.AppProcess.MemoryUsed,
	)
}

func (sh *SystemHealth) CheckProcess() {
	p, err := findBackend(8080)
	if err != nil {
		log.Println("Failed to fetch process info: ", err)
		sh.AppProcess.IsRunning = false
		return
	}
	sh.AppProcess.IsRunning = true

	c, err := p.CPUPercent()
	if err != nil {
		log.Println("Error reading process CPU usage: ", err)
		return
	} else {
		sh.AppProcess.CpuPercent = c
	}

	m, err := p.MemoryInfo()
	if err != nil {
		log.Println("Error reading process CPU usage: ", err)
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

func (sh *SystemHealth) Listen() {
	sh.UpdateStatus()
	sh.CheckProcess()
	ticker := time.NewTicker(sh.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			sh.UpdateStatus()
			sh.CheckProcess()
			fmt.Println(sh.StatusTextSummery())
		case newInterval := <-sh.IntervalChannel:
			ticker.Reset(newInterval)
			sh.UpdateStatus()
		}
	}
}

func (sh *SystemHealth) UpdateStatus() {
	// RAM Stats
	v, err := mem.VirtualMemory()
	if err != nil {
		log.Println("Error reading Virtual Memory stats: ", err)
	}

	sh.Status.Memory.Update(
		v.Used,
		v.Total,
		v.SwapCached,
		v.SwapTotal,
	)

	// CPU Stats
	c, err := cpu.Percent(time.Second, false)
	if err != nil {
		log.Println("Error reading CPU stats: ", err)
	} else {
		sh.Status.CPU.Update(c[0])
	}

	// Disk Stats
	d, err := disk.Usage("/")
	if err != nil {
		log.Println("Error reading Disk stats: ", err)
	} else {
		sh.Status.Disk.Update(d.Used, d.Total)
	}

	// Uptime
	u, err := host.Uptime()
	if err != nil {
		log.Println("Error reading Uptime: ", err)
	} else {
		sh.Status.Uptime = u
	}

}
