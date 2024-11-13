package winprocess

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/shirou/gopsutil/v4/process"
)

// ProcessInfo represents information about a process.
type ProcessInfo struct {
	PID        int32     `json:"pid"`
	Name       string    `json:"name"`
	Executable string    `json:"executable"`
	CPUPercent float64   `json:"cpu_percent"`
	Memory     float32   `json:"memory_percent"`
	Status     string    `json:"status"`
	CreateTime time.Time `json:"create_time"`
}

// GetProcesses returns a list of all running processes with their information.
func GetProcesses() ([]*ProcessInfo, error) {
	processes, err := process.Processes()
	if err != nil {
		return nil, fmt.Errorf("failed to get processes: %w", err)
	}

	var processList []*ProcessInfo
	for _, p := range processes {
		// Skip processes that cause errors
		if p.Pid == 4 { // Or check for other known problematic PIDs
			continue
		}

		info, err := getProcessInfo(p)
		if err != nil {
			// Log the error but continue with other processes
			fmt.Printf("Failed to get process info for PID %d: %v\n", p.Pid, err)
			continue
		}
		processList = append(processList, info)
	}

	return processList, nil
}

// getProcessInfo retrieves information about a single process.
func getProcessInfo(p *process.Process) (*ProcessInfo, error) {
	name, err := p.Name()
	if err != nil {
		return nil, fmt.Errorf("failed to get process name: %w", err)
	}

	executable, err := p.Exe()
	if err != nil {
		return nil, fmt.Errorf("failed to get process executable: %w", err)
	}

	cpuPercent, err := p.CPUPercent()
	if err != nil {
		return nil, fmt.Errorf("failed to get CPU percent: %w", err)
	}

	memPercent, err := p.MemoryPercent()
	if err != nil {
		return nil, fmt.Errorf("failed to get memory percent: %w", err)
	}

	statusSlice, err := p.Status()
	var status string
	if err != nil {
		// Handle the error gracefully
		fmt.Printf("Failed to get process status for PID %d: %v\n", p.Pid, err)
		status = "unknown" // Or provide a more informative message
	} else {
		status = statusSlice[0] // Assuming the first status is the primary one
	}

	createTime, err := p.CreateTime()
	if err != nil {
		return nil, fmt.Errorf("failed to get process create time: %w", err)
	}

	return &ProcessInfo{
		PID:        p.Pid,
		Name:       name,
		Executable: executable,
		CPUPercent: cpuPercent,
		Memory:     memPercent,
		Status:     status,
		CreateTime: time.Unix(0, createTime*int64(time.Millisecond)),
	}, nil
}

// KillProcess terminates a process with the given PID.
func KillProcess(pid int32) error {
	p, err := process.NewProcess(pid)
	if err != nil {
		return fmt.Errorf("failed to find process with PID %d: %w", pid, err)
	}

	err = p.Kill()
	if err != nil {
		return fmt.Errorf("failed to kill process: %w", err)
	}

	return nil
}

// RestartProcess restarts a process with the given PID using its executable path and arguments.
func RestartProcess(pid int32) error {
	p, err := process.NewProcess(pid)
	if err != nil {
		return fmt.Errorf("failed to find process with PID %d: %w", pid, err)
	}

	executable, err := p.Exe()
	if err != nil {
		return fmt.Errorf("failed to get process executable: %w", err)
	}

	// Get the command-line arguments
	cmdline, err := p.CmdlineSlice()
	if err != nil {
		return fmt.Errorf("failed to get process arguments: %w", err)
	}

	err = p.Kill()
	if err != nil {
		return fmt.Errorf("failed to kill process: %w", err)
	}

	// Construct the command with arguments
	cmd := exec.Command(executable, cmdline[1:]...)
	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to restart process: %w", err)
	}

	return nil
}
