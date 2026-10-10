package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// Theme Colors - Tokyo Night / Modern Cyberpunk
var (
	ColorPrimary   = lipgloss.Color("#7AA2F7") // Vibrant Blue
	ColorSecondary = lipgloss.Color("#BB9AF7") // Neon Purple
	ColorSuccess   = lipgloss.Color("#9ECE6A") // Mint Green
	ColorWarning   = lipgloss.Color("#E0AF68") // Warm Amber
	ColorDanger    = lipgloss.Color("#F7768E") // Coral Red
	ColorCyan      = lipgloss.Color("#7DCFFF") // Bright Cyan
	ColorMuted     = lipgloss.Color("#565F89") // Muted Slate
	ColorBgDark    = lipgloss.Color("#1A1B26") // Deep Background
	ColorCardBg    = lipgloss.Color("#24283B") // Card Background
	ColorHighlight = lipgloss.Color("#414868") // Active Row
	ColorText      = lipgloss.Color("#C0CAF5") // Foreground Text
	ColorWhite     = lipgloss.Color("#FFFFFF")
)

// Global Styles
var (
	// Header & Logo
	LogoStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorCyan).
			Background(ColorCardBg).
			Padding(0, 1)

	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(ColorMuted)

	// Tab Styles
	ActiveTabStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorWhite).
			Background(ColorPrimary).
			Padding(0, 2)

	InactiveTabStyle = lipgloss.NewStyle().
				Foreground(ColorMuted).
				Background(ColorCardBg).
				Padding(0, 2)

	TabGapStyle = lipgloss.NewStyle().
			Foreground(ColorMuted).
			Padding(0, 1)

	// Box Containers
	CardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorPrimary).
			Padding(1, 1).
			Background(ColorBgDark)

	BorderMuted = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorMuted).
			Padding(0, 1)

	// Table Rows
	HeaderRowStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorSecondary).
			Background(ColorCardBg).
			Padding(0, 1)

	SelectedRowStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(ColorWhite).
				Background(ColorHighlight).
				Padding(0, 1)

	NormalRowStyle = lipgloss.NewStyle().
			Foreground(ColorText).
			Padding(0, 1)

	// Badges
	BadgeTCP = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorCyan).
			Background(ColorCardBg).
			Padding(0, 1)

	BadgeUDP = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorWarning).
			Background(ColorCardBg).
			Padding(0, 1)

	BadgeSuccess = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorSuccess)

	BadgeDanger = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorDanger)

	// Status & Toast
	ToastStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorWhite).
			Background(ColorSecondary).
			Padding(0, 1)

	KeyHelpStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorCyan)

	DescHelpStyle = lipgloss.NewStyle().
			Foreground(ColorMuted)
)
