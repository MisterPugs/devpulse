package model

import (
	"fmt"
	"sort"
	"strings"

	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
)

// PortInfo represents a single open port entry.
type PortInfo struct {
	Port        uint32
	Protocol    string
	IP          string
	PID         int32
	ProcessName string
	MemoryMB    float64
	User        string
	Status      string
}

// FetchOpenPorts collects all active listening ports with associated process information.
func FetchOpenPorts() ([]PortInfo, error) {
	conns, err := net.Connections("all")
	if err != nil {
		return nil, fmt.Errorf("unable to scan connections: %w", err)
	}

	// Cache process info by PID to avoid redundant calls
	procMap := make(map[int32]*process.Process)
	nameMap := make(map[int32]string)
	memMap := make(map[int32]float64)

	var results []PortInfo
	seen := make(map[string]bool)

	for _, c := range conns {
		isListening := c.Status == "LISTEN" || (c.Type == 2 && c.Status == "NONE") // TCP LISTEN or UDP
		if !isListening || c.Laddr.Port == 0 {
			continue
		}

		proto := "TCP"
		if c.Type == 2 {
			proto = "UDP"
		}

		key := fmt.Sprintf("%s:%d:%d", proto, c.Laddr.Port, c.Pid)
		if seen[key] {
			continue
		}
		seen[key] = true

		procName := "Unknown"
		var memMB float64

		if c.Pid > 0 {
			if n, ok := nameMap[c.Pid]; ok {
				procName = n
				memMB = memMap[c.Pid]
			} else {
				p, ok := procMap[c.Pid]
				if !ok {
					var pErr error
					p, pErr = process.NewProcess(c.Pid)
					if pErr == nil {
						procMap[c.Pid] = p
					}
				}

				if p != nil {
					if n, err := p.Name(); err == nil && n != "" {
						procName = n
					}
					if memInfo, err := p.MemoryInfo(); err == nil && memInfo != nil {
						memMB = float64(memInfo.RSS) / 1024 / 1024
					}
					nameMap[c.Pid] = procName
					memMap[c.Pid] = memMB
				}
			}
		}

		ipStr := c.Laddr.IP
		if ipStr == "" || ipStr == "0.0.0.0" || ipStr == "::" {
			ipStr = "*"
		}

		results = append(results, PortInfo{
			Port:        c.Laddr.Port,
			Protocol:    proto,
			IP:          ipStr,
			PID:         c.Pid,
			ProcessName: procName,
			MemoryMB:    memMB,
			Status:      c.Status,
		})
	}

	// Sort primarily by port number
	sort.Slice(results, func(i, j int) bool {
		return results[i].Port < results[j].Port
	})

	return results, nil
}

// FilterPorts filters a list of PortInfo based on a search query.
func FilterPorts(ports []PortInfo, query string) []PortInfo {
	if strings.TrimSpace(query) == "" {
		return ports
	}
	q := strings.ToLower(query)
	var filtered []PortInfo
	for _, p := range ports {
		portStr := fmt.Sprintf("%d", p.Port)
		pidStr := fmt.Sprintf("%d", p.PID)
		if strings.Contains(portStr, q) ||
			strings.Contains(strings.ToLower(p.ProcessName), q) ||
			strings.Contains(strings.ToLower(p.Protocol), q) ||
			strings.Contains(strings.ToLower(p.IP), q) ||
			strings.Contains(pidStr, q) {
			filtered = append(filtered, p)
		}
	}
	return filtered
}
