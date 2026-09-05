package commands

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/Awak3r/PortKiller/internal/port"
)

// Filter selects processes by name and/or port. Construct with
// NewFilter; a zero value is not usable. A nil port means "not set".
type Filter struct {
	name      string
	port      *int
	collector Collector
	killer    Killer
}

// Option customizes Filter dependencies.
type Option func(*Filter)

// WithCollector overrides the process source (test seam).
func WithCollector(c Collector) Option { return func(f *Filter) { f.collector = c } }

// WithKiller overrides the terminator (test seam).
func WithKiller(k Killer) Option { return func(f *Filter) { f.killer = k } }

// Collector provides the snapshot of listening processes.
type Collector interface {
	Collect() ([]port.ProcessInfo, error)
}

// Killer terminates a process by pid.
type Killer interface {
	KillByPid(pid int32) error
}

type prodCollector struct{}

func (prodCollector) Collect() ([]port.ProcessInfo, error) { return port.Collect() }

type prodKiller struct{}

func (prodKiller) KillByPid(pid int32) error { return port.KillByPid(pid) }

// NewFilter validates the port when it is provided and applies options.
func NewFilter(name string, port *int, opts ...Option) (Filter, error) {
	if port != nil {
		if err := ValidatePort(*port); err != nil {
			return Filter{}, err
		}
	}
	f := Filter{name: name, port: port, collector: prodCollector{}, killer: prodKiller{}}
	for _, opt := range opts {
		opt(&f)
	}
	return f, nil
}

// List prints matching processes into w (nil -> os.Stdout). No matches
// is not an error: an empty table is printed.
func (f Filter) List(w io.Writer) error {
	procs, err := f.selectProcesses()
	if err != nil {
		return err
	}
	printTable(w, procs)
	return nil
}

// Kill terminates matching processes with a worker pool. Returns
// found/killed counts over actually dispatched processes; partial
// failures are joined into one error, the stats are still reported.
func (f Filter) Kill(ctx context.Context) (int, int, error) {
	procs, err := f.selectProcesses()
	if err != nil {
		return 0, 0, err
	}
	if len(procs) == 0 {
		return 0, 0, ErrProcessNotFound
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

	sent := 0
	for _, proc := range procs {
		if proc.Pid <= 0 {
			mu.Lock()
			errs = append(errs, fmt.Errorf("cannot kill process on port %d: owned by another user, rerun under sudo", proc.Port))
			mu.Unlock()
			continue
		}
		select {
		case jobs <- proc:
			sent++
		case <-ctx.Done():
			close(jobs)
			wg.Wait()
			return sent, killed, errors.Join(append(errs, ctx.Err())...)
		}
	}
	close(jobs)
	wg.Wait()

	return sent, killed, errors.Join(errs...)
}

func (f Filter) match(p ProcessInfo) bool {
	if f.name != "" && !NameMatches(p.Name, f.name) {
		return false
	}
	if f.port != nil && p.Port != *f.port {
		return false
	}
	return true
}

func (f Filter) selectProcesses() ([]ProcessInfo, error) {
	procs, err := f.collector.Collect()
	if err != nil {
		return nil, err
	}
	res := []ProcessInfo{}
	for _, proc := range procs {
		if f.match(proc) {
			res = append(res, proc)
		}
	}
	return res, nil
}
