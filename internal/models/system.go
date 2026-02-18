package models

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
