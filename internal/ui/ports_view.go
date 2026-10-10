package ui

import (
	"fmt"
	"strings"

	"github.com/MisterPugs/devpulse/internal/model"
	"github.com/charmbracelet/lipgloss"
)

func renderPortsView(ports []model.PortInfo, cursor int, filter string, width int, height int) string {
	var b strings.Builder

	// Search bar line
	searchBar := ""
	if filter != "" {
		searchBar = fmt.Sprintf(" Filter: [%s] (press 'esc' to clear)", filter)
	} else {
		searchBar = " Press '/' to filter ports (e.g. 3000, node, python)"
	}
	b.WriteString(SubtitleStyle.Render(searchBar) + "\n\n")

	if len(ports) == 0 {
		emptyMsg := "No listening ports found matching the filter."
		if filter == "" {
			emptyMsg = "Scanning active listening ports..."
		}
		b.WriteString(lipgloss.NewStyle().Foreground(ColorWarning).Render(emptyMsg) + "\n")
		return b.String()
	}

	// Table header
	header := fmt.Sprintf("%-7s  %-5s  %-24s  %-8s  %-12s  %-10s",
		"PORT", "PROTO", "PROCESS NAME", "PID", "BIND IP", "MEMORY")
	b.WriteString(HeaderRowStyle.Width(width - 4).Render(header) + "\n")

	// Calculate visible window for scrolling
	visibleRows := height - 10
	if visibleRows < 5 {
		visibleRows = 5
	}

	start := 0
	if cursor >= visibleRows {
		start = cursor - visibleRows + 1
	}
	end := start + visibleRows
	if end > len(ports) {
		end = len(ports)
	}

	for i := start; i < end; i++ {
		p := ports[i]
		protoBadge := "TCP"
		if p.Protocol == "UDP" {
			protoBadge = "UDP"
		}

		memStr := "-"
		if p.MemoryMB > 0 {
			memStr = fmt.Sprintf("%.1f MB", p.MemoryMB)
		}

		procName := p.ProcessName
		if len(procName) > 24 {
			procName = procName[:21] + "..."
		}

		line := fmt.Sprintf(":%-6d %-5s  %-24s  %-8d  %-12s  %-10s",
			p.Port, protoBadge, procName, p.PID, p.IP, memStr)

		if i == cursor {
			b.WriteString(SelectedRowStyle.Width(width - 4).Render("▶ " + line) + "\n")
		} else {
			b.WriteString(NormalRowStyle.Width(width - 4).Render("  " + line) + "\n")
		}
	}

	// Footer counter
	counter := fmt.Sprintf("\n Showing %d of %d active listening ports (Use 'k' to Kill, 'r' to Refresh)", len(ports), len(ports))
	b.WriteString(SubtitleStyle.Render(counter))

	return b.String()
}
