package commands

import (
	"errors"
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"github.com/Awak3r/PortKiller/internal/port"
)

var ErrProcessNotFound = errors.New("process not found")

type ProcessInfo = port.ProcessInfo

const killWorkers = 4

func (f Filter) match(p ProcessInfo) bool {
	if f.Name != "" && !NameMatches(p.Name, f.Name) {
		return false
	}
	if f.Port != nil && p.Port != *f.Port {
		return false
	}
	return true
}

func (f Filter) selectProcesses() ([]ProcessInfo, error) {
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
