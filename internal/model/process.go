package model

import (
	"fmt"
	"sort"
	"strings"

	"github.com/shirou/gopsutil/v3/process"
)

// ProcessInfo represents a running process summary.
type ProcessInfo struct {
	PID      int32
	Name     string
	CPU      float64
	MemoryMB float64
	MemPct   float32
	Status   string
	Cmdline  string
}

// SortBy specifies the sorting criteria for processes.
type SortBy int

const (
	SortByCPU SortBy = iota
	SortByMemory
	SortByName
	SortByPID
)

// FetchProcesses retrieves a list of top processes.
func FetchProcesses(limit int) ([]ProcessInfo, error) {
	pids, err := process.Pids()
	if err != nil {
		return nil, fmt.Errorf("could not list process IDs: %w", err)
	}

	var list []ProcessInfo
	for _, pid := range pids {
		if pid <= 0 {
			continue
		}
		p, err := process.NewProcess(pid)
		if err != nil {
			continue
		}

		name, err := p.Name()
		if err != nil || name == "" {
			name = "System"
		}

		cpu, _ := p.CPUPercent()
		memInfo, _ := p.MemoryInfo()
		memPct, _ := p.MemoryPercent()

		var memMB float64
		if memInfo != nil {
			memMB = float64(memInfo.RSS) / (1024 * 1024)
		}

		status, _ := p.Status()
		statusStr := "R"
		if len(status) > 0 {
			statusStr = status[0]
		}

		cmd, _ := p.Cmdline()
		if len(cmd) > 60 {
			cmd = cmd[:57] + "..."
		}

		list = append(list, ProcessInfo{
			PID:      pid,
			Name:     name,
			CPU:      cpu,
			MemoryMB: memMB,
			MemPct:   memPct,
			Status:   statusStr,
			Cmdline:  cmd,
		})
	}

	return list, nil
}

// SortProcesses sorts the process slice in-place by the specified column.
func SortProcesses(procs []ProcessInfo, sortBy SortBy) {
	switch sortBy {
	case SortByCPU:
		sort.Slice(procs, func(i, j int) bool {
			return procs[i].CPU > procs[j].CPU
		})
	case SortByMemory:
		sort.Slice(procs, func(i, j int) bool {
			return procs[i].MemoryMB > procs[j].MemoryMB
		})
	case SortByName:
		sort.Slice(procs, func(i, j int) bool {
			return strings.ToLower(procs[i].Name) < strings.ToLower(procs[j].Name)
		})
	case SortByPID:
		sort.Slice(procs, func(i, j int) bool {
			return procs[i].PID < procs[j].PID
		})
	}
}

// FilterProcesses filters the process slice by keyword.
func FilterProcesses(procs []ProcessInfo, query string) []ProcessInfo {
	if strings.TrimSpace(query) == "" {
		return procs
	}
	q := strings.ToLower(query)
	var filtered []ProcessInfo
	for _, p := range procs {
		pidStr := fmt.Sprintf("%d", p.PID)
		if strings.Contains(strings.ToLower(p.Name), q) ||
			strings.Contains(strings.ToLower(p.Cmdline), q) ||
			strings.Contains(pidStr, q) {
			filtered = append(filtered, p)
		}
	}
	return filtered
}
