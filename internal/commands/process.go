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
