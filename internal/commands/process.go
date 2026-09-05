// Package commands implements the actions behind the portkiller CLI:
// filtering listening processes (list) and terminating them (kill).
//
// Both commands share a single pipeline: Collect -> filter -> action.
// Name matching is case-insensitive substring (NameMatches);
// port is a 1-65535 integer validated before any sudo escalation.
package commands

import (
	"errors"
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"github.com/Awak3r/PortKiller/internal/port"
)

// ErrProcessNotFound is returned when no listening process matches the filter.
var ErrProcessNotFound = errors.New("process not found")

// ProcessInfo re-exported for callers of this package.
type ProcessInfo = port.ProcessInfo

// Collect proxies the port package so callers keep a single import.
func Collect() ([]ProcessInfo, error) { return port.Collect() }

// procFilter selects processes by name and/or port.
// Zero values mean "no constraint" for the corresponding field.
type procFilter struct {
	name string
	port int
}

// newFilter validates the port up front and builds a filter.
func newFilter(name string, portNum int, portSet bool) (procFilter, error) {
	if portSet {
		if portNum < 1 || portNum > 65535 {
			return procFilter{}, fmt.Errorf("invalid port (1-65535)")
		}
		return procFilter{name: name, port: portNum}, nil
	}
	return procFilter{name: name}, nil
}

// match reports whether the process satisfies the filter
// (NameMatches semantics — review item 4).
func (f procFilter) match(p ProcessInfo) bool {
	if f.name != "" && !NameMatches(p.Name, f.name) {
		return false
	}
	if f.port != 0 && p.Port != f.port {
		return false
	}
	return true
}

// select runs the collection pipeline shared by list and kill.
// Returns ErrProcessNotFound (nil slice) when nothing matches.
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

// printTable renders processes as an aligned PORT | PID | NAME table
// into w (or os.Stdout when w is nil).
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

// List renders the table into w (or os.Stdout when w is nil) — the writer
// is injected so tests can capture the output (review item 8).
// List validates the port when explicitly set, filters listening processes
// and prints them as a table into w (nil -> os.Stdout).
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

// Kill filters listening processes and terminates each of them.
// Returns found/killed counts; partial failures are joined into one error,
// but the stats are still reported (review item 6).
func Kill(filter procFilter) (int, int, error) {
	procs, err := filter.selectProcesses()
	if err != nil {
		return 0, 0, err
	}
	killed := 0
	var errs []error
	for _, proc := range procs {
		if err := port.KillByPid(int32(proc.Pid)); err != nil {
			errs = append(errs, fmt.Errorf("failed to kill process %s (PID: %d): %w", proc.Name, proc.Pid, err))
		} else {
			killed++
		}
	}
	return len(procs), killed, errors.Join(errs...)
}
