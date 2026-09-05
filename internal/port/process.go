package port

import (
	"errors"
	"fmt"
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

// Время на «нежное» завершение перед SIGKILL.
const (
	killGracePeriod = 2 * time.Second
	killPollPeriod  = 100 * time.Millisecond
)

// KillByPid завершает процесс в два этапа:
//
//  1. SIGTERM — «вежливая» просьба завершиться: процесс может перехватить
//     сигнал и корректно закрыть сокеты, сбросить буферы, удалить pidfile.
//  2. Если после killGracePeriod процесс всё ещё жив (завис, игнорирует
//     SIGTERM) — отправляется SIGKILL. Его нельзя перехватить: ядро
//     убивает процесс немедленно.
//
// Если процесс исчез между снимком списка и отправкой сигнала (ESRCH),
// ошибка не возвращается — цель уже достигнута.
func KillByPid(pid int32) error {
	proc, err := process.NewProcess(pid)
	if err != nil {
		return fmt.Errorf("invalid pid %d: %w", pid, err)
	}

	// этап 1: мягкое завершение
	err = proc.Terminate()
	if err != nil {
		if errors.Is(err, syscall.ESRCH) {
			return nil // процесс уже завершился сам
		}
		return fmt.Errorf("failed to send SIGTERM to process %d: %w", pid, err)
	}

	// grace period: ждём завершения, опрашивая процесс
	deadline := time.Now().Add(killGracePeriod)
	for time.Now().Before(deadline) {
		time.Sleep(killPollPeriod)
		if alive, err := proc.IsRunning(); err != nil || !alive {
			return nil // завершился (или статус недоступен — считаем завершённым)
		}
	}

	// этап 2: жёсткое убийство
	if err := proc.Kill(); err != nil && !errors.Is(err, syscall.ESRCH) {
		return fmt.Errorf("process %d ignored SIGTERM and SIGKILL failed: %w", pid, err)
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
