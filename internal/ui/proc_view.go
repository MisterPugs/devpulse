package ui

import (
	"fmt"
	"strings"

	"github.com/MisterPugs/devpulse/internal/model"
	"github.com/charmbracelet/lipgloss"
)

func renderProcView(procs []model.ProcessInfo, cursor int, sortBy model.SortBy, filter string, width int, height int) string {
	var b strings.Builder

	// Sort & Filter bar
	sortName := "CPU"
	switch sortBy {
	case model.SortByMemory:
		sortName = "Memory"
	case model.SortByName:
		sortName = "Name"
	case model.SortByPID:
		sortName = "PID"
	}

	sortBar := fmt.Sprintf(" Sorted by: [%s] (press 'c' CPU, 'm' Mem, 'n' Name, 'p' PID)", sortName)
	if filter != "" {
		sortBar += fmt.Sprintf(" | Filter: [%s]", filter)
	}
	b.WriteString(SubtitleStyle.Render(sortBar) + "\n\n")

	if len(procs) == 0 {
		b.WriteString(lipgloss.NewStyle().Foreground(ColorWarning).Render("No processes found matching filter.") + "\n")
		return b.String()
	}

	header := fmt.Sprintf("%-8s  %-24s  %-8s  %-8s  %-12s  %-30s",
		"PID", "PROCESS NAME", "CPU %", "MEM %", "RSS", "COMMAND")
	b.WriteString(HeaderRowStyle.Width(width - 4).Render(header) + "\n")

	visibleRows := height - 10
	if visibleRows < 5 {
		visibleRows = 5
	}

	start := 0
	if cursor >= visibleRows {
		start = cursor - visibleRows + 1
	}
	end := start + visibleRows
	if end > len(procs) {
		end = len(procs)
	}

	for i := start; i < end; i++ {
		p := procs[i]
		name := p.Name
		if len(name) > 24 {
			name = name[:21] + "..."
		}

		cpuStr := fmt.Sprintf("%.1f%%", p.CPU)
		memPctStr := fmt.Sprintf("%.1f%%", p.MemPct)
		rssStr := fmt.Sprintf("%.1f MB", p.MemoryMB)

		cmd := p.Cmdline
		if len(cmd) > 30 {
			cmd = cmd[:27] + "..."
		}

		line := fmt.Sprintf("%-8d  %-24s  %-8s  %-8s  %-12s  %-30s",
			p.PID, name, cpuStr, memPctStr, rssStr, cmd)

		if i == cursor {
			b.WriteString(SelectedRowStyle.Width(width - 4).Render("▶ " + line) + "\n")
		} else {
			b.WriteString(NormalRowStyle.Width(width - 4).Render("  " + line) + "\n")
		}
	}

	counter := fmt.Sprintf("\n Showing %d processes (Press 'k' to Kill selected process)", len(procs))
	b.WriteString(SubtitleStyle.Render(counter))

	return b.String()
}
