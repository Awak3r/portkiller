package utils

import (
	"fmt"
	"log"
	"os"
	"os/exec"
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

func Collect() []ProcessInfo {
	EnsureRoot()
	processes := []ProcessInfo{}
	seen := make(map[struct{ pid, port int }]struct{})
	conns, err := net.Connections("tcp4")
	if err != nil {
		log.Fatal(err)
	}
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
		k := struct{ pid, port int }{pid, port}
		if _, dup := seen[k]; dup {
			continue
		}
		seen[k] = struct{}{}
		cur := ProcessInfo{Name: name, Pid: pid, Port: port}
		processes = append(processes, cur)
	}
	sort.Slice(processes, func(i, j int) bool { return processes[i].Port < processes[j].Port })
	return processes
}

func EnsureRoot() {
	if os.Geteuid() == 0 {
		return
	}
	exe, err := os.Executable()
	if err != nil {
		os.Exit(1)
	}
	cmd := exec.Command("sudo", append([]string{exe}, os.Args[1:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
