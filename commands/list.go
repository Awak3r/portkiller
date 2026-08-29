package commands

import (
	"fmt"
	"github.com/Awak3r/PortKiller/process"
	"os"
	"strings"
	"text/tabwriter"
)

func printTable(p []process.ProcessInfo) {
	w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintln(w, "PROCESS\tPORT\tPID")
	fmt.Fprintln(w, "-------\t----\t----")
	for _, proc := range p {
		fmt.Fprintf(w, "%s\t%d\t%d\n", proc.Name, proc.Port, proc.Pid)
	}
	w.Flush()
}

func ListAll(p []process.ProcessInfo) {
	printTable(p)
}

func ListByName(name string, p []process.ProcessInfo) {
	res := []process.ProcessInfo{}
	name = strings.ToLower(name)
	for _, proc := range p {
		if strings.HasPrefix(strings.ToLower(proc.Name), name) {
			res = append(res, proc)
		}
	}
	printTable(res)
}

func ListByPort(port int, p []process.ProcessInfo) {
	res := []process.ProcessInfo{}
	for _, proc := range p {
		if port == proc.Port {
			res = append(res, proc)
		}
	}
	printTable(res)
}
