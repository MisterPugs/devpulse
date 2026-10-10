package ui

import (
	"fmt"
	"strings"

	"github.com/MisterPugs/devpulse/internal/model"
	"github.com/charmbracelet/lipgloss"
)

func renderSysView(snap *model.SystemSnapshot, width int) string {
	if snap == nil {
		return lipgloss.NewStyle().Foreground(ColorWarning).Render("Gathering system performance metrics...")
	}

	var b strings.Builder

	// Top System Banner
	sysHeader := fmt.Sprintf(" Host: %s | OS: %s (%s) | Uptime: %s",
		snap.Hostname, snap.OS, snap.Platform, snap.Uptime.Round(1e9))
	b.WriteString(lipgloss.NewStyle().Bold(true).Foreground(ColorCyan).Render(sysHeader) + "\n\n")

	// CPU Section
	b.WriteString(TitleStyle.Render("⚡ CPU Activity & Real-time Sparkline") + "\n")
	cpuBar := model.ProgressBar(snap.CPUPercent, 35)
	spark := model.Sparkline(snap.CPUHistory, 100)

	b.WriteString(fmt.Sprintf(" Load: %s\n", cpuBar))
	if spark != "" {
		b.WriteString(fmt.Sprintf(" Trend: %s\n", lipgloss.NewStyle().Foreground(ColorSecondary).Render(spark)))
	}
	b.WriteString("\n")

	// Memory Section
	b.WriteString(TitleStyle.Render("🧠 Memory Utilization (RAM & Swap)") + "\n")
	memBar := model.ProgressBar(snap.MemUsedPct, 35)
	memDetails := fmt.Sprintf(" Used: %.2f GB / Total: %.2f GB (Free: %.2f GB)",
		snap.MemUsedGB, snap.MemTotalGB, snap.MemFreeGB)
	b.WriteString(fmt.Sprintf(" RAM:  %s  %s\n", memBar, SubtitleStyle.Render(memDetails)))

	if snap.SwapTotalGB > 0 {
		swapBar := model.ProgressBar(snap.SwapUsedPct, 35)
		swapDetails := fmt.Sprintf(" Used: %.2f GB / Total: %.2f GB",
			snap.SwapUsedGB, snap.SwapTotalGB)
		b.WriteString(fmt.Sprintf(" Swap: %s  %s\n", swapBar, SubtitleStyle.Render(swapDetails)))
	}
	b.WriteString("\n")

	// Network I/O
	b.WriteString(TitleStyle.Render("🌐 Network Throughput") + "\n")
	netStats := fmt.Sprintf(" ↓ Download: %.2f KB/s    ↑ Upload: %.2f KB/s    (Total In: %.1f MB / Out: %.1f MB)",
		snap.NetDnloadSec, snap.NetUploadSec,
		float64(snap.NetBytesRecv)/1024/1024,
		float64(snap.NetBytesSent)/1024/1024)
	b.WriteString(lipgloss.NewStyle().Foreground(ColorSuccess).Render(netStats) + "\n\n")

	// Disks
	if len(snap.Disks) > 0 {
		b.WriteString(TitleStyle.Render("💾 Storage & Mounted Volumes") + "\n")
		for _, d := range snap.Disks {
			dBar := model.ProgressBar(d.UsedPct, 25)
			dDetails := fmt.Sprintf(" %.1f GB / %.1f GB (Free: %.1f GB) [%s]",
				d.UsedGB, d.TotalGB, d.FreeGB, d.FSType)
			b.WriteString(fmt.Sprintf(" %-10s %s %s\n", d.MountPoint, dBar, SubtitleStyle.Render(dDetails)))
		}
	}

	return b.String()
}
