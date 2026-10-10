package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/MisterPugs/devpulse/internal/killer"
	"github.com/MisterPugs/devpulse/internal/model"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(1500*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// AppModel represents the full application state.
type AppModel struct {
	activeTab int
	width     int
	height    int

	// Data
	ports       []model.PortInfo
	filteredPorts []model.PortInfo
	procs       []model.ProcessInfo
	filteredProcs []model.ProcessInfo
	sysSnap     *model.SystemSnapshot
	cleanItems  []model.CleanItem
	scanningClean bool

	// Cursors
	portsCursor int
	procsCursor int
	cleanCursor int

	// Search & Filter
	filterMode  bool
	filterQuery string
	sortBy      model.SortBy

	// Toast & Modal
	toastMsg  string
	toastTimer time.Time
	showHelp  bool
}

// NewAppModel initializes the Bubble Tea model.
func NewAppModel() AppModel {
	return AppModel{
		activeTab:   0,
		sortBy:      model.SortByCPU,
		width:       80,
		height:      24,
	}
}

func (m AppModel) Init() tea.Cmd {
	return tea.Batch(
		m.refreshDataCmd(),
		tickCmd(),
	)
}

func (m AppModel) refreshDataCmd() tea.Cmd {
	return func() tea.Msg {
		ports, _ := model.FetchOpenPorts()
		procs, _ := model.FetchProcesses(100)
		sys, _ := model.FetchSystemSnapshot()
		return dataLoadedMsg{
			ports: ports,
			procs: procs,
			sys:   sys,
		}
	}
}

type dataLoadedMsg struct {
	ports []model.PortInfo
	procs []model.ProcessInfo
	sys   *model.SystemSnapshot
}

type cleanLoadedMsg struct {
	items []model.CleanItem
}

func (m AppModel) scanCleanCmd() tea.Cmd {
	return func() tea.Msg {
		items, _ := model.ScanDevBloat(".", 4)
		return cleanLoadedMsg{items: items}
	}
}

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tickMsg:
		cmds = append(cmds, m.refreshDataCmd(), tickCmd())
		if time.Since(m.toastTimer) > 4*time.Second {
			m.toastMsg = ""
		}

	case dataLoadedMsg:
		m.ports = msg.ports
		m.filteredPorts = model.FilterPorts(m.ports, m.filterQuery)
		if m.portsCursor >= len(m.filteredPorts) && len(m.filteredPorts) > 0 {
			m.portsCursor = len(m.filteredPorts) - 1
		}

		m.procs = msg.procs
		model.SortProcesses(m.procs, m.sortBy)
		m.filteredProcs = model.FilterProcesses(m.procs, m.filterQuery)
		if m.procsCursor >= len(m.filteredProcs) && len(m.filteredProcs) > 0 {
			m.procsCursor = len(m.filteredProcs) - 1
		}

		m.sysSnap = msg.sys

	case cleanLoadedMsg:
		m.cleanItems = msg.items
		m.scanningClean = false

	case tea.KeyMsg:
		if m.filterMode {
			switch msg.String() {
			case "enter", "esc":
				m.filterMode = false
			case "backspace":
				if len(m.filterQuery) > 0 {
					m.filterQuery = m.filterQuery[:len(m.filterQuery)-1]
					m.applyFilter()
				}
			default:
				if len(msg.String()) == 1 {
					m.filterQuery += msg.String()
					m.applyFilter()
				}
			}
			return m, nil
		}

		// Global keys
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "?":
			m.showHelp = !m.showHelp
			return m, nil

		case "esc":
			if m.showHelp {
				m.showHelp = false
				return m, nil
			}
			if m.filterQuery != "" {
				m.filterQuery = ""
				m.applyFilter()
				return m, nil
			}

		case "/":
			m.filterMode = true
			return m, nil

		case "tab":
			m.activeTab = (m.activeTab + 1) % 4
			m.filterQuery = ""
			m.applyFilter()
			if m.activeTab == 3 && len(m.cleanItems) == 0 && !m.scanningClean {
				m.scanningClean = true
				cmds = append(cmds, m.scanCleanCmd())
			}

		case "shift+tab":
			m.activeTab = (m.activeTab + 3) % 4
			m.filterQuery = ""
			m.applyFilter()

		case "1":
			m.activeTab = 0
		case "2":
			m.activeTab = 1
		case "3":
			m.activeTab = 2
		case "4":
			m.activeTab = 3
			if len(m.cleanItems) == 0 && !m.scanningClean {
				m.scanningClean = true
				cmds = append(cmds, m.scanCleanCmd())
			}

		case "r":
			m.setToast("Refreshing data...")
			if m.activeTab == 3 {
				m.scanningClean = true
				cmds = append(cmds, m.scanCleanCmd())
			} else {
				cmds = append(cmds, m.refreshDataCmd())
			}

		// Navigation & Actions per Tab
		case "up":
			switch m.activeTab {
			case 0:
				if m.portsCursor > 0 {
					m.portsCursor--
				}
			case 1:
				if m.procsCursor > 0 {
					m.procsCursor--
				}
			case 3:
				if m.cleanCursor > 0 {
					m.cleanCursor--
				}
			}

		case "down":
			switch m.activeTab {
			case 0:
				if m.portsCursor < len(m.filteredPorts)-1 {
					m.portsCursor++
				}
			case 1:
				if m.procsCursor < len(m.filteredProcs)-1 {
					m.procsCursor++
				}
			case 3:
				if m.cleanCursor < len(m.cleanItems)-1 {
					m.cleanCursor++
				}
			}

		case "c":
			if m.activeTab == 1 {
				m.sortBy = model.SortByCPU
				model.SortProcesses(m.filteredProcs, m.sortBy)
				m.setToast("Sorted by CPU Usage")
			}

		case "m":
			if m.activeTab == 1 {
				m.sortBy = model.SortByMemory
				model.SortProcesses(m.filteredProcs, m.sortBy)
				m.setToast("Sorted by Memory")
			}

		case "n":
			if m.activeTab == 1 {
				m.sortBy = model.SortByName
				model.SortProcesses(m.filteredProcs, m.sortBy)
				m.setToast("Sorted by Process Name")
			}

		case "p":
			if m.activeTab == 1 {
				m.sortBy = model.SortByPID
				model.SortProcesses(m.filteredProcs, m.sortBy)
				m.setToast("Sorted by PID")
			}

		case "k", "x", "delete":
			// Kill action
			if m.activeTab == 0 && len(m.filteredPorts) > 0 {
				target := m.filteredPorts[m.portsCursor]
				if target.PID > 0 {
					err := killer.KillPID(target.PID, true)
					if err != nil {
						m.setToast(fmt.Sprintf("Failed to kill PID %d: %v", target.PID, err))
					} else {
						m.setToast(fmt.Sprintf("Terminated PID %d (%s) on port :%d", target.PID, target.ProcessName, target.Port))
						cmds = append(cmds, m.refreshDataCmd())
					}
				}
			} else if m.activeTab == 1 && len(m.filteredProcs) > 0 {
				target := m.filteredProcs[m.procsCursor]
				err := killer.KillPID(target.PID, true)
				if err != nil {
					m.setToast(fmt.Sprintf("Failed to kill %s (PID %d): %v", target.Name, target.PID, err))
				} else {
					m.setToast(fmt.Sprintf("Terminated %s (PID %d)", target.Name, target.PID))
					cmds = append(cmds, m.refreshDataCmd())
				}
			}

		case " ":
			// Toggle clean item
			if m.activeTab == 3 && len(m.cleanItems) > 0 {
				m.cleanItems[m.cleanCursor].Selected = !m.cleanItems[m.cleanCursor].Selected
			}

		case "a":
			// Select all clean items
			if m.activeTab == 3 && len(m.cleanItems) > 0 {
				allSel := true
				for _, it := range m.cleanItems {
					if !it.Selected {
						allSel = false
						break
					}
				}
				for i := range m.cleanItems {
					m.cleanItems[i].Selected = !allSel
				}
			}

		case "d":
			// Delete selected clean items
			if m.activeTab == 3 {
				var deletedCount int
				var freedBytes int64
				var remaining []model.CleanItem
				for _, item := range m.cleanItems {
					if item.Selected {
						if err := model.DeleteItem(item); err == nil {
							deletedCount++
							freedBytes += item.SizeByte
						} else {
							remaining = append(remaining, item)
						}
					} else {
						remaining = append(remaining, item)
					}
				}
				m.cleanItems = remaining
				m.cleanCursor = 0
				m.setToast(fmt.Sprintf("Cleaned %d folders, reclaimed %s", deletedCount, model.FormatBytes(freedBytes)))
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *AppModel) applyFilter() {
	m.filteredPorts = model.FilterPorts(m.ports, m.filterQuery)
	m.portsCursor = 0
	m.filteredProcs = model.FilterProcesses(m.procs, m.filterQuery)
	m.procsCursor = 0
}

func (m *AppModel) setToast(msg string) {
	m.toastMsg = msg
	m.toastTimer = time.Now()
}

func (m AppModel) View() string {
	var b strings.Builder

	// Header line
	logo := LogoStyle.Render("⚡ DEVPULSE")
	tagline := SubtitleStyle.Render(" The Modern Developer Cockpit & Process TUI")
	headerLeft := lipgloss.JoinHorizontal(lipgloss.Center, logo, " ", tagline)

	var toast string
	if m.toastMsg != "" {
		toast = ToastStyle.Render(" " + m.toastMsg + " ")
	}
	header := lipgloss.JoinHorizontal(lipgloss.Top, headerLeft, "  ", toast)
	b.WriteString(header + "\n\n")

	// Navigation Tabs
	tabs := []string{"[ 1 ] 🔌 Ports", "[ 2 ] 🚀 Processes", "[ 3 ] 📊 System Vitals", "[ 4 ] 🧹 Dev Cleaner"}
	var renderedTabs []string
	for i, t := range tabs {
		if i == m.activeTab {
			renderedTabs = append(renderedTabs, ActiveTabStyle.Render(t))
		} else {
			renderedTabs = append(renderedTabs, InactiveTabStyle.Render(t))
		}
	}
	tabRow := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	b.WriteString(tabRow + "\n\n")

	// Main Content Area inside Card
	var content string
	switch m.activeTab {
	case 0:
		content = renderPortsView(m.filteredPorts, m.portsCursor, m.filterQuery, m.width, m.height)
	case 1:
		content = renderProcView(m.filteredProcs, m.procsCursor, m.sortBy, m.filterQuery, m.width, m.height)
	case 2:
		content = renderSysView(m.sysSnap, m.width)
	case 3:
		content = renderCleanView(m.cleanItems, m.cleanCursor, m.scanningClean, m.width, m.height)
	}

	card := CardStyle.Width(m.width - 4).Height(m.height - 8).Render(content)
	b.WriteString(card + "\n")

	// Footer Help / Shortcuts
	footer := KeyHelpStyle.Render(" [Tab]") + DescHelpStyle.Render(" Switch Tab  ") +
		KeyHelpStyle.Render("[↑/↓]") + DescHelpStyle.Render(" Navigate  ") +
		KeyHelpStyle.Render("[/]") + DescHelpStyle.Render(" Filter  ") +
		KeyHelpStyle.Render("[k]") + DescHelpStyle.Render(" Kill  ") +
		KeyHelpStyle.Render("[r]") + DescHelpStyle.Render(" Refresh  ") +
		KeyHelpStyle.Render("[?]") + DescHelpStyle.Render(" Help  ") +
		KeyHelpStyle.Render("[q]") + DescHelpStyle.Render(" Quit")
	b.WriteString(footer)

	// Modal Overlay if help requested
	if m.showHelp {
		modal := renderHelpModal(m.width)
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, modal)
	}

	return b.String()
}
