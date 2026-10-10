package killer

import (
	"fmt"
	"os/exec"
	"runtime"
	"strconv"

	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
)

// KillPID terminates a process by its PID across Windows, Linux, and macOS.
func KillPID(pid int32, force bool) error {
	p, err := process.NewProcess(pid)
	if err != nil {
		return fmt.Errorf("process with PID %d not found: %w", pid, err)
	}

	name, _ := p.Name()

	// Try graceful termination first if not force
	if runtime.GOOS == "windows" {
		flag := "/T"
		if force {
			flag = "/F"
		}
		cmd := exec.Command("taskkill", flag, "/PID", strconv.Itoa(int(pid)))
		if err := cmd.Run(); err != nil {
			// Fallback to gopsutil Kill
			if kErr := p.Kill(); kErr != nil {
				return fmt.Errorf("failed to kill PID %d (%s): %w", pid, name, err)
			}
		}
		return nil
	}

	// Unix / macOS
	if force {
		return p.Kill() // SIGKILL
	}
	return p.Terminate() // SIGTERM
}

// FindPIDByPort finds the PID listening on a given port.
func FindPIDByPort(targetPort uint32) (int32, string, error) {
	conns, err := net.Connections("all")
	if err != nil {
		return 0, "", fmt.Errorf("failed to retrieve network connections: %w", err)
	}

	for _, c := range conns {
		if c.Laddr.Port == targetPort && (c.Status == "LISTEN" || c.Status == "NONE" || c.Status == "") {
			if c.Pid > 0 {
				procName := "unknown"
				if p, pErr := process.NewProcess(c.Pid); pErr == nil {
					if n, nErr := p.Name(); nErr == nil {
						procName = n
					}
				}
				return c.Pid, procName, nil
			}
		}
	}

	return 0, "", fmt.Errorf("no process found listening on port %d", targetPort)
}

// KillPort finds and kills the process listening on a given port.
func KillPort(targetPort uint32, force bool) (int32, string, error) {
	pid, name, err := FindPIDByPort(targetPort)
	if err != nil {
		return 0, "", err
	}

	if err := KillPID(pid, force); err != nil {
		return pid, name, err
	}

	return pid, name, nil
}
