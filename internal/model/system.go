package model

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

// DiskInfo represents filesystem stats.
type DiskInfo struct {
	MountPoint string
	FSType     string
	TotalGB    float64
	UsedGB     float64
	FreeGB     float64
	UsedPct    float64
}

// SystemSnapshot holds real-time system performance metrics.
type SystemSnapshot struct {
	Hostname     string
	OS           string
	Platform     string
	Uptime       time.Duration
	CPUPercent   float64
	PerCPU       []float64
	CPUHistory   []float64
	MemTotalGB   float64
	MemUsedGB    float64
	MemFreeGB    float64
	MemUsedPct   float64
	SwapTotalGB  float64
	SwapUsedGB   float64
	SwapUsedPct  float64
	Disks        []DiskInfo
	NetBytesSent uint64
	NetBytesRecv uint64
	NetUploadSec float64 // KB/s
	NetDnloadSec float64 // KB/s
}

var (
	lastNetSent uint64
	lastNetRecv uint64
	lastNetTime time.Time
	cpuHistory  []float64
)

const maxHistory = 40

// FetchSystemSnapshot captures current hardware performance and vitals.
func FetchSystemSnapshot() (*SystemSnapshot, error) {
	snap := &SystemSnapshot{}

	// Host info
	if hInfo, err := host.Info(); err == nil && hInfo != nil {
		snap.Hostname = hInfo.Hostname
		snap.OS = hInfo.OS
		snap.Platform = hInfo.Platform
		snap.Uptime = time.Duration(hInfo.Uptime) * time.Second
	}

	// CPU info
	percentages, err := cpu.Percent(0, false)
	if err == nil && len(percentages) > 0 {
		snap.CPUPercent = percentages[0]
	}

	// Update CPU history for sparklines
	cpuHistory = append(cpuHistory, snap.CPUPercent)
	if len(cpuHistory) > maxHistory {
		cpuHistory = cpuHistory[len(cpuHistory)-maxHistory:]
	}
	snap.CPUHistory = make([]float64, len(cpuHistory))
	copy(snap.CPUHistory, cpuHistory)

	// Per-core CPU
	if per, err := cpu.Percent(0, true); err == nil {
		snap.PerCPU = per
	}

	// Virtual Memory
	if vm, err := mem.VirtualMemory(); err == nil && vm != nil {
		snap.MemTotalGB = float64(vm.Total) / (1024 * 1024 * 1024)
		snap.MemUsedGB = float64(vm.Used) / (1024 * 1024 * 1024)
		snap.MemFreeGB = float64(vm.Available) / (1024 * 1024 * 1024)
		snap.MemUsedPct = vm.UsedPercent
	}

	// Swap Memory
	if sm, err := mem.SwapMemory(); err == nil && sm != nil {
		snap.SwapTotalGB = float64(sm.Total) / (1024 * 1024 * 1024)
		snap.SwapUsedGB = float64(sm.Used) / (1024 * 1024 * 1024)
		snap.SwapUsedPct = sm.UsedPercent
	}

	// Disks
	if parts, err := disk.Partitions(false); err == nil {
		for _, p := range parts {
			if usage, err := disk.Usage(p.Mountpoint); err == nil && usage != nil && usage.Total > 0 {
				snap.Disks = append(snap.Disks, DiskInfo{
					MountPoint: p.Mountpoint,
					FSType:     p.Fstype,
					TotalGB:    float64(usage.Total) / (1024 * 1024 * 1024),
					UsedGB:     float64(usage.Used) / (1024 * 1024 * 1024),
					FreeGB:     float64(usage.Free) / (1024 * 1024 * 1024),
					UsedPct:    usage.UsedPercent,
				})
			}
		}
	}

	// Network I/O
	if ioCounters, err := net.IOCounters(false); err == nil && len(ioCounters) > 0 {
		var totalSent, totalRecv uint64
		for _, io := range ioCounters {
			totalSent += io.BytesSent
			totalRecv += io.BytesRecv
		}

		now := time.Now()
		if !lastNetTime.IsZero() {
			dt := now.Sub(lastNetTime).Seconds()
			if dt > 0 {
				if totalSent >= lastNetSent {
					snap.NetUploadSec = float64(totalSent-lastNetSent) / 1024 / dt
				}
				if totalRecv >= lastNetRecv {
					snap.NetDnloadSec = float64(totalRecv-lastNetRecv) / 1024 / dt
				}
			}
		}

		lastNetSent = totalSent
		lastNetRecv = totalRecv
		lastNetTime = now
		snap.NetBytesSent = totalSent
		snap.NetBytesRecv = totalRecv
	}

	return snap, nil
}

// Sparkline renders a slice of floats as a Unicode sparkline ( ▂▃▄▅▆▇█).
func Sparkline(data []float64, maxVal float64) string {
	if len(data) == 0 {
		return ""
	}
	bars := []rune{' ', ' ', '▂', '▃', '▄', '▅', '▆', '▇', '█'}
	if maxVal <= 0 {
		maxVal = 100
	}

	var runes []rune
	for _, val := range data {
		if val < 0 {
			val = 0
		}
		if val > maxVal {
			val = maxVal
		}
		idx := int((val / maxVal) * float64(len(bars)-1))
		if idx >= len(bars) {
			idx = len(bars) - 1
		}
		runes = append(runes, bars[idx])
	}
	return string(runes)
}

// ProgressBar renders a colored or text-based progress bar.
func ProgressBar(pct float64, width int) string {
	if width <= 2 {
		return ""
	}
	if pct < 0 {
		pct = 0
	}
	if pct > 100 {
		pct = 100
	}
	fillCount := int((pct / 100.0) * float64(width))
	if fillCount > width {
		fillCount = width
	}

	filled := ""
	for i := 0; i < fillCount; i++ {
		filled += "█"
	}
	empty := ""
	for i := fillCount; i < width; i++ {
		empty += "░"
	}
	return fmt.Sprintf("[%s%s] %5.1f%%", filled, empty, pct)
}
