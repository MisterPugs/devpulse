package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func renderHelpModal(width int) string {
	var b strings.Builder

	b.WriteString(TitleStyle.Render("⚡ DEVPULSE - KEYBOARD SHORTCUTS") + "\n\n")

	keys := []struct {
		Key  string
		Desc string
	}{
		{"1, 2, 3, 4", "Switch Tabs directly (Ports / Procs / System / Clean)"},
		{"Tab / Shift+Tab", "Next / Previous Tab"},
		{"↑ / k, ↓ / j", "Navigate rows up / down"},
		{"/ (slash)", "Start real-time search / filter"},
		{"Esc", "Clear filter or cancel action"},
		{"k / x", "Kill selected process or port (Ports & Procs tabs)"},
		{"r", "Refresh data immediately"},
		{"c, m, n, p", "Sort processes by CPU, Memory, Name, or PID"},
		{"Space", "Toggle selection (Dev Cleaner tab)"},
		{"a", "Select all items (Dev Cleaner tab)"},
		{"d", "Delete selected artifacts (Dev Cleaner tab)"},
		{"?", "Toggle this Help modal"},
		{"q / Ctrl+C", "Quit DevPulse"},
	}

	for _, k := range keys {
		line := lipgloss.JoinHorizontal(
			lipgloss.Left,
			KeyHelpStyle.Width(18).Render(k.Key),
			DescHelpStyle.Render(k.Desc),
		)
		b.WriteString(line + "\n")
	}

	b.WriteString("\n" + SubtitleStyle.Render("Press '?' or 'Esc' to close this modal"))

	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(ColorSecondary).
		Background(ColorCardBg).
		Padding(1, 2).
		Width(65)

	return modalStyle.Render(b.String())
}
