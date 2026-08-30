package commands

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Awak3r/PortKiller/process"
)

var ErrProcessNotFound = errors.New("process not found")

func KillByName(name string) (int, int, error) {
	p, err := process.Collect()
	if err != nil {
		return 0, 0, err
	}
	found := 0
	killed := 0
	var errs []error
	for _, proc := range p {
		if strings.EqualFold(proc.Name, name) {
			found++
			err := process.KillByPid(int32(proc.Pid))
			if err != nil {
				errs = append(errs, fmt.Errorf("failed to kill process %s (PID: %d): %w", proc.Name, proc.Pid, err))
			} else {
				killed++
			}
		}
	}
	if found == 0 {
		return 0, 0, ErrProcessNotFound
	}
	return found, killed, errors.Join(errs...)
}

func KillByPort(port int) (int, int, error) {
	if port < 1 || port > 65535 {
		return 0, 0, fmt.Errorf("invalid port (1-65535)")
	}
	p, err := process.Collect()
	if err != nil {
		return 0, 0, err
	}
	found := 0
	killed := 0
	var errs []error
	for _, proc := range p {
		if port == proc.Port {
			found++
			err := process.KillByPid(int32(proc.Pid))
			if err != nil {
				errs = append(errs, fmt.Errorf("failed to kill process %d (PID: %d): %w", proc.Port, proc.Pid, err))
			} else {
				killed++
			}
		}
	}
	if found == 0 {
		return 0, 0, ErrProcessNotFound
	}
	return found, killed, errors.Join(errs...)
}

func KillByNameAndPort(name string, port int) (int, int, error) {
	if port < 1 || port > 65535 {
		return 0, 0, fmt.Errorf("invalid port (1-65535)")
	}
	p, err := process.Collect()
	if err != nil {
		return 0, 0, err
	}
	found := 0
	killed := 0
	var errs []error
	for _, proc := range p {
		if strings.EqualFold(proc.Name, name) && port == proc.Port {
			found++
			err := process.KillByPid(int32(proc.Pid))
			if err != nil {
				errs = append(errs, fmt.Errorf("failed to kill process %s (PID: %d): %w", proc.Name, proc.Pid, err))
			} else {
				killed++
			}
		}
	}
	if found == 0 {
		return 0, 0, ErrProcessNotFound
	}
	return found, killed, errors.Join(errs...)
}
