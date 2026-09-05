package commands

import (
	"fmt"
	"github.com/Awak3r/PortKiller/internal/port"
	"os"
	"text/tabwriter"
)

func printTable(p []port.ProcessInfo) {
	w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintln(w, "PROCESS\tPORT\tPID")
	fmt.Fprintln(w, "-------\t----\t----")
	for _, proc := range p {
		fmt.Fprintf(w, "%s\t%d\t%d\n", proc.Name, proc.Port, proc.Pid)
	}
	w.Flush()
}

func ListAll() error {
	p, err := port.Collect()
	if err != nil {
		return err
	}
	printTable(p)
	return nil
}

func ListByName(name string) error {
	p, err := port.Collect()
	if err != nil {
		return err
	}
	res := []port.ProcessInfo{}
	for _, proc := range p {
		if NameMatches(proc.Name, name) {
			res = append(res, proc)
		}
	}
	printTable(res)
	return nil
}

func ListByPort(portNum int) error {
	if portNum < 1 || portNum > 65535 {
		return fmt.Errorf("invalid port (1-65535)")
	}
	p, err := port.Collect()
	if err != nil {
		return err
	}
	res := []port.ProcessInfo{}
	for _, proc := range p {
		if portNum == proc.Port {
			res = append(res, proc)
		}
	}
	printTable(res)
	return nil
}

func ListByNameAndPort(name string, portNum int) error {
	if portNum < 1 || portNum > 65535 {
		return fmt.Errorf("invalid port (1-65535)")
	}
	p, err := port.Collect()
	if err != nil {
		return err
	}
	res := []port.ProcessInfo{}
	for _, proc := range p {
		if NameMatches(proc.Name, name) && portNum == proc.Port {
			res = append(res, proc)
		}
	}
	printTable(res)
	return nil
}
