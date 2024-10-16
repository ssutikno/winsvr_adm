package process

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v4/process"
)

// ProcessInfo represents information about a process.
type ProcessInfo struct {
	PID        int32     `json:"pid"`
	Name       string    `json:"name"`
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

	cpuPercent, err := p.CPUPercent()
	if err != nil {
		return nil, fmt.Errorf("failed to get CPU percent: %w", err)
	}

	memPercent, err := p.MemoryPercent()
	if err != nil {
		return nil, fmt.Errorf("failed to get memory percent: %w", err)
	}

	status, err := p.Status()
	if err != nil {
		return nil, fmt.Errorf("failed to get process status: %w", err)
	}

	createTime, err := p.CreateTime()
	if err != nil {
		return nil, fmt.Errorf("failed to get process create time: %w", err)
	}

	return &ProcessInfo{
		PID:        p.Pid,
		Name:       name,
		CPUPercent: cpuPercent,
		Memory:     memPercent,
		Status:     status,
		CreateTime: time.Unix(0, createTime*int64(time.Millisecond)), // Convert to time.Time
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

// RestartProcess restarts a process with the given PID.
// Note: This function currently just kills and then tries to start the process again
// using its name. A more robust implementation might involve using process managers
// or service supervisors.
func RestartProcess(pid int32) error {
	p, err := process.NewProcess(pid)
	if err != nil {
		return fmt.Errorf("failed to find process with PID %d: %w", pid, err)
	}

	name, err := p.Name()
	if err != nil {
		return fmt.Errorf("failed to get process name: %w", err)
	}

	err = p.Kill()
	if err != nil {
		return fmt.Errorf("failed to kill process: %w", err)
	}

	// This is a simplified restart. Replace with your actual restart logic.
	_, err = process.NewProcess(0).Cmd().Output() // Replace with actual command to start the process
	if err != nil {
		return fmt.Errorf("failed to restart process: %w", err)
	}

	return nil
}
