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
type Filter struct {
	Name string
	Port *int
}

// NewFilter validates the port when it is provided.
func NewFilter(name string, port *int) (Filter, error) {
	if port != nil {
		if err := ValidatePort(*port); err != nil {
			return Filter{}, err
		}
	}
	return Filter{Name: name, Port: port}, nil
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
