package port

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"syscall"
	"time"

	"github.com/shirou/gopsutil/v4/net"
	"github.com/shirou/gopsutil/v4/process"
)

type ProcessInfo struct {
	Name string
	Pid  int
	Port int
}

const (
	killGracePeriod = 2 * time.Second
	killPollPeriod  = 100 * time.Millisecond
)

// KillByPid terminates a process in two stages: SIGTERM, then SIGKILL
// after a grace period if the process is still alive. A process that
// vanished (ESRCH) or became a zombie counts as terminated.
func KillByPid(pid int32) error {
	proc, err := process.NewProcess(pid)
	if err != nil {
		return fmt.Errorf("invalid pid %d: %w", pid, err)
	}

	err = proc.Terminate()
	if err != nil {
		if isProcessGone(err) {
			return nil
		}
		return fmt.Errorf("failed to send SIGTERM to process %d: %w", pid, err)
	}

	deadline := time.Now().Add(killGracePeriod)
	for time.Now().Before(deadline) {
		time.Sleep(killPollPeriod)
		alive, err := isAlive(proc)
		if err != nil || !alive {
			return nil
		}
	}

	if err := proc.Kill(); err != nil && !isProcessGone(err) {
		return fmt.Errorf("process %d ignored SIGTERM and SIGKILL failed: %w", pid, err)
	}
	return nil
}

// isProcessGone reports whether the error means the process no longer
// exists. os.ErrProcessDone is what os.Process.Signal returns for an
// already-finished process; ESRCH covers direct syscall paths.
func isProcessGone(err error) bool {
	return errors.Is(err, os.ErrProcessDone) || errors.Is(err, syscall.ESRCH)
}

// isAlive reports whether the process is still a running (non-zombie,
// existing) process. Errors are surfaced so a transient /proc failure
// is treated as "still alive" (conservative: leads to SIGKILL, never
// to a false success report).
func isAlive(proc *process.Process) (bool, error) {
	statuses, err := proc.Status()
	if err != nil {
		return false, err
	}
	for _, s := range statuses {
		if s == "Z" {
			return false, nil
		}
	}
	return true, nil
}

func Collect() ([]ProcessInfo, error) {
	conns, err := net.Connections("tcp")
	if err != nil {
		return nil, fmt.Errorf("failed to get network connections: %w", err)
	}
	seen := make(map[string]struct{})
	var processes []ProcessInfo
	for _, conn := range conns {
		if conn.Status != "LISTEN" {
			continue
		}
		var name string
		var pid int
		if conn.Pid > 0 {
			p, err := process.NewProcess(conn.Pid)
			if err != nil {
				continue
			}
			name, err = p.Name()
			if err != nil {
				name = "unknown"
			}
			pid = int(conn.Pid)
		} else {
			name = "-"
			pid = 0
		}
		port := int(conn.Laddr.Port)
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
	sort.SliceStable(processes, func(i, j int) bool {
		if processes[i].Port != processes[j].Port {
			return processes[i].Port < processes[j].Port
		}
		return processes[i].Pid < processes[j].Pid
	})
	return processes, nil
}
