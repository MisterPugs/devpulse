package ui

import (
	"fmt"
	"strings"

	"github.com/MisterPugs/devpulse/internal/model"
	"github.com/charmbracelet/lipgloss"
)

func renderCleanView(items []model.CleanItem, cursor int, scanning bool, width int, height int) string {
	var b strings.Builder

	b.WriteString(SubtitleStyle.Render(" Scan current workspace for bloated build artifacts and heavy dependencies.") + "\n")
	b.WriteString(SubtitleStyle.Render(" [Space] Toggle Select | [a] Select All | [d] Delete Selected | [r] Rescan") + "\n\n")

	if scanning {
		b.WriteString(lipgloss.NewStyle().Foreground(ColorWarning).Render("⏳ Scanning directory tree for dev artifacts...") + "\n")
		return b.String()
	}

	if len(items) == 0 {
		b.WriteString(lipgloss.NewStyle().Foreground(ColorSuccess).Render("✨ Workspace is clean! No heavy dev caches found.") + "\n")
		return b.String()
	}

	var totalReclaimable int64
	var selectedReclaimable int64
	for _, item := range items {
		totalReclaimable += item.SizeByte
		if item.Selected {
			selectedReclaimable += item.SizeByte
		}
	}

	summary := fmt.Sprintf(" Total Artifacts: %s | Selected to Clean: %s",
		model.FormatBytes(totalReclaimable),
		lipgloss.NewStyle().Bold(true).Foreground(ColorDanger).Render(model.FormatBytes(selectedReclaimable)))
	b.WriteString(summary + "\n\n")

	header := fmt.Sprintf("%-4s  %-28s  %-12s  %-40s",
		"SEL", "CATEGORY", "SIZE", "DIRECTORY PATH")
	b.WriteString(HeaderRowStyle.Width(width - 4).Render(header) + "\n")

	visibleRows := height - 12
	if visibleRows < 5 {
		visibleRows = 5
	}

	start := 0
	if cursor >= visibleRows {
		start = cursor - visibleRows + 1
	}
	end := start + visibleRows
	if end > len(items) {
		end = len(items)
	}

	for i := start; i < end; i++ {
		item := items[i]
		checkbox := "[ ]"
		if item.Selected {
			checkbox = lipgloss.NewStyle().Foreground(ColorDanger).Render("[✓]")
		}

		pathStr := item.Path
		if len(pathStr) > 40 {
			pathStr = "..." + pathStr[len(pathStr)-37:]
		}

		line := fmt.Sprintf("%-4s  %-28s  %-12s  %-40s",
			checkbox, item.Category, model.FormatBytes(item.SizeByte), pathStr)

		if i == cursor {
			b.WriteString(SelectedRowStyle.Width(width - 4).Render("▶ " + line) + "\n")
		} else {
			b.WriteString(NormalRowStyle.Width(width - 4).Render("  " + line) + "\n")
		}
	}

	return b.String()
}
