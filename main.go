package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/MisterPugs/devpulse/internal/killer"
	"github.com/MisterPugs/devpulse/internal/model"
	"github.com/MisterPugs/devpulse/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

const version = "1.0.0"

func printBanner() {
	fmt.Printf("\n⚡ DevPulse v%s — Fast Developer Cockpit & Process TUI\n", version)
}

func printHelp() {
	printBanner()
	fmt.Println(`
USAGE:
  devpulse                  Launch the interactive TUI dashboard
  devpulse ports            Quickly list listening ports to terminal
  devpulse kill <port>      Find & kill process listening on specified port
  devpulse clean            Scan current workspace for dev artifacts & bloat
  devpulse --version, -v    Show version
  devpulse --help, -h       Show this help message

EXAMPLES:
  devpulse
  devpulse kill 3000
  devpulse kill :8080
  devpulse ports
  devpulse clean
`)
}

func main() {
	if len(os.Args) > 1 {
		arg := strings.ToLower(os.Args[1])
		switch arg {
		case "--help", "-h", "help":
			printHelp()
			return

		case "--version", "-v", "version":
			fmt.Printf("devpulse version %s\n", version)
			return

		case "ports":
			printBanner()
			ports, err := model.FetchOpenPorts()
			if err != nil {
				fmt.Printf("Error scanning ports: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("\n%-7s  %-5s  %-24s  %-8s  %-12s\n", "PORT", "PROTO", "PROCESS", "PID", "BIND IP")
			fmt.Println(strings.Repeat("─", 65))
			for _, p := range ports {
				fmt.Printf(":%-6d %-5s  %-24s  %-8d  %-12s\n", p.Port, p.Protocol, p.ProcessName, p.PID, p.IP)
			}
			fmt.Printf("\nTotal: %d listening ports.\n", len(ports))
			return

		case "kill":
			if len(os.Args) < 3 {
				fmt.Println("Usage: devpulse kill <port> (e.g. devpulse kill 3000)")
				os.Exit(1)
			}
			portRaw := strings.TrimPrefix(os.Args[2], ":")
			portNum, err := strconv.ParseUint(portRaw, 10, 32)
			if err != nil {
				fmt.Printf("Invalid port number '%s'\n", os.Args[2])
				os.Exit(1)
			}
			pid, name, err := killer.KillPort(uint32(portNum), true)
			if err != nil {
				fmt.Printf("Failed to kill port :%d: %v\n", portNum, err)
				os.Exit(1)
			}
			fmt.Printf("✓ Successfully terminated %s (PID %d) on port :%d\n", name, pid, portNum)
			return

		case "clean":
			printBanner()
			fmt.Println("Scanning current workspace for build artifacts...")
			items, err := model.ScanDevBloat(".", 4)
			if err != nil {
				fmt.Printf("Error scanning: %v\n", err)
				os.Exit(1)
			}
			if len(items) == 0 {
				fmt.Println("✓ No heavy dev artifacts found. Workspace is clean!")
				return
			}
			var total int64
			fmt.Printf("\n%-24s  %-12s  %s\n", "CATEGORY", "SIZE", "PATH")
			fmt.Println(strings.Repeat("─", 70))
			for _, item := range items {
				total += item.SizeByte
				fmt.Printf("%-24s  %-12s  %s\n", item.Category, model.FormatBytes(item.SizeByte), item.Path)
			}
			fmt.Printf("\nTotal reclaimable: %s across %d folders.\n", model.FormatBytes(total), len(items))
			return
		}
	}

	// Interactive TUI mode
	p := tea.NewProgram(ui.NewAppModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running DevPulse: %v\n", err)
		os.Exit(1)
	}
}
