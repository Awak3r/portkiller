package commands

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
	"text/tabwriter"

	"github.com/Awak3r/PortKiller/internal/port"
)

var ErrProcessNotFound = errors.New("process not found")

type ProcessInfo = port.ProcessInfo

func Collect() ([]ProcessInfo, error) { return port.Collect() }

type procFilter struct {
	name string
	port int
}

func newFilter(name string, portNum int, portSet bool) (procFilter, error) {
	if portSet {
		if err := validatePort(portNum); err != nil {
			return procFilter{}, err
		}
		return procFilter{name: name, port: portNum}, nil
	}
	return procFilter{name: name}, nil
}

func (f procFilter) match(p ProcessInfo) bool {
	if f.name != "" && !NameMatches(p.Name, f.name) {
		return false
	}
	if f.port != 0 && p.Port != f.port {
		return false
	}
	return true
}

func (f procFilter) selectProcesses() ([]ProcessInfo, error) {
	procs, err := port.Collect()
	if err != nil {
		return nil, err
	}
	res := []ProcessInfo{}
	for _, proc := range procs {
		if f.match(proc) {
			res = append(res, proc)
		}
	}
	if len(res) == 0 {
		return nil, ErrProcessNotFound
	}
	return res, nil
}

func printTable(w io.Writer, procs []ProcessInfo) {
	if w == nil {
		w = os.Stdout
	}
	tab := tabwriter.NewWriter(w, 0, 8, 2, ' ', 0)
	fmt.Fprintln(tab, "PROCESS\tPORT\tPID")
	fmt.Fprintln(tab, "-------\t----\t----")
	for _, proc := range procs {
		fmt.Fprintf(tab, "%s\t%d\t%d\n", proc.Name, proc.Port, proc.Pid)
	}
	tab.Flush()
}

func List(w io.Writer, name string, portNum int, portSet bool) error {
	filter, err := newFilter(name, portNum, portSet)
	if err != nil {
		return err
	}
	procs, err := filter.selectProcesses()
	if err != nil {
		return err
	}
	printTable(w, procs)
	return nil
}

const killWorkers = 4

func Kill(ctx context.Context, filter procFilter) (int, int, error) {
	procs, err := filter.selectProcesses()
	if err != nil {
		return 0, 0, err
	}

	jobs := make(chan ProcessInfo)
	var wg sync.WaitGroup
	var mu sync.Mutex
	killed := 0
	var errs []error

	for i := 0; i < killWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for proc := range jobs {
				if err := port.KillByPid(int32(proc.Pid)); err != nil {
					mu.Lock()
					errs = append(errs, fmt.Errorf("failed to kill process %s (PID: %d): %w", proc.Name, proc.Pid, err))
					mu.Unlock()
				} else {
					mu.Lock()
					killed++
					mu.Unlock()
				}
			}
		}()
	}

	for _, proc := range procs {
		select {
		case jobs <- proc:
		case <-ctx.Done():
			close(jobs)
			wg.Wait()
			return len(procs), killed, errors.Join(append(errs, ctx.Err())...)
		}
	}
	close(jobs)
	wg.Wait()

	return len(procs), killed, errors.Join(errs...)
}
