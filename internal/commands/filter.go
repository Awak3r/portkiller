package commands

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/Awak3r/PortKiller/internal/port"
)

// Filter selects processes by name and/or port.
// A nil Port means "not set"; zero values of Name mean "not set".
// Collector and Killer are seams for tests; in production they are
// backed by the port package.
type Filter struct {
	Name string
	Port *int

	collector Collector
	killer    Killer
}

// Collector provides the snapshot of listening processes.
type Collector interface {
	Collect() ([]port.ProcessInfo, error)
}

// Killer terminates a process by pid.
type Killer interface {
	KillByPid(pid int32) error
}

// collector is the production Collector.
type prodCollector struct{}

func (prodCollector) Collect() ([]port.ProcessInfo, error) { return port.Collect() }

// prodKiller is the production Killer.
type prodKiller struct{}

func (prodKiller) KillByPid(pid int32) error { return port.KillByPid(pid) }

// NewFilter validates the port when it is provided.
func NewFilter(name string, port *int) (Filter, error) {
	if port != nil {
		if err := ValidatePort(*port); err != nil {
			return Filter{}, err
		}
	}
	f := Filter{Name: name, Port: port, collector: prodCollector{}, killer: prodKiller{}}
	return f, nil
}

// List validates the filter and prints matching processes into w
// (nil -> os.Stdout).
func (f Filter) List(w io.Writer) error {
	procs, err := f.selectProcesses()
	if err != nil {
		return err
	}
	printTable(w, procs)
	return nil
}

// Kill terminates matching processes with a worker pool.
// Returns found/killed counts; partial failures are joined into one
// error, but the stats are still reported.
func (f Filter) Kill(ctx context.Context) (int, int, error) {
	procs, err := f.selectProcesses()
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
				if err := f.killer.KillByPid(int32(proc.Pid)); err != nil {
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
		if proc.Pid <= 0 {
			mu.Lock()
			errs = append(errs, fmt.Errorf("cannot kill process on port %d: owned by another user, rerun under sudo", proc.Port))
			mu.Unlock()
			continue
		}
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
