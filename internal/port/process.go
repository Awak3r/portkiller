package port

import (
	"fmt"
	"sort"

	"github.com/shirou/gopsutil/v4/net"
	"github.com/shirou/gopsutil/v4/process"
)

type ProcessInfo struct {
	Name string
	Pid  int
	Port int
}

func KillByPid(pid int32) error {
	proc, err := process.NewProcess(pid)
	if err != nil {
		return fmt.Errorf("invalid pid %d: %w", pid, err)
	}
	err = proc.Terminate()
	if err != nil {
		return fmt.Errorf("failed to kill process %d: %w", pid, err)
	}
	return nil
}

func Collect() ([]ProcessInfo, error) {
	conns, err := net.Connections("tcp4")
	if err != nil {
		return nil, fmt.Errorf("failed to get network connections: %w", err)
	}
	seen := make(map[string]struct{})
	var processes []ProcessInfo
	for _, conn := range conns {
		if conn.Status != "LISTEN" {
			continue
		}
		p, err := process.NewProcess(conn.Pid)
		if err != nil {
			continue
		}
		name, err := p.Name()
		if err != nil {
			name = "unknown"
		}
		port := int(conn.Laddr.Port)
		pid := int(conn.Pid)
		key := fmt.Sprintf("%d-%d", pid, port)
		if _, dup := seen[key]; dup {
			continue
		}
		seen[key] = struct{}{}
		processes = append(processes, ProcessInfo{
			Name: name,
			Pid:  pid,
			Port: port,
		})
	}
	sort.Slice(processes, func(i, j int) bool {
		return processes[i].Port < processes[j].Port
	})
	return processes, nil
}
