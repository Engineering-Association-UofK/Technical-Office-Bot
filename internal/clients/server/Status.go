package server

import (
	"log/slog"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
)

type CPU struct {
	Load    float64
	CoreNum int
}

func (c *CPU) Update(load float64) {
	c.Load = load
}

type Memory struct {
	Current     uint64
	Max         uint64
	CurrentSwap uint64
	MaxSwap     uint64
}

func (m *Memory) Update(current uint64, max uint64, currentSwap uint64, maxSwap uint64) {
	m.Current = current
	m.Max = max
	m.CurrentSwap = currentSwap
	m.MaxSwap = maxSwap
}

type Disk struct {
	Current uint64
	Max     uint64
}

func (s *Disk) Update(current uint64, max uint64) {
	s.Current = current
	s.Max = max
}

type Status struct {
	Memory Memory
	Disk   Disk
	CPU    CPU

	Uptime uint64
}

type AppProcess struct {
	IsRunning  bool
	CpuPercent float64
	MemoryUsed uint64
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
