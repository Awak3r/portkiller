package commands

import "fmt"

// PrintUsage is kept for compatibility with the legacy help entry point.
func PrintUsage() {
	fmt.Println(`portkiller — kill processes by name or port

Usage:
  portkiller list [-name NAME | -port PORT]
  portkiller kill [-name NAME | -port PORT]
  portkiller -version
  portkiller -help`)
}

// ---- list entry points ----

func ListAll() error {
	return List(procFilter{})
}

func ListByName(name string) error {
	return List(procFilter{name: name})
}

func ListByPort(portNum int) error {
	f, err := newFilter("", portNum, true)
	if err != nil {
		return err
	}
	return List(f)
}

func ListByNameAndPort(name string, portNum int) error {
	f, err := newFilter(name, portNum, true)
	if err != nil {
		return err
	}
	return List(f)
}

// ---- kill entry points ----

func KillByName(name string) (int, int, error) {
	return Kill(procFilter{name: name})
}

func KillByPort(portNum int) (int, int, error) {
	f, err := newFilter("", portNum, true)
	if err != nil {
		return 0, 0, err
	}
	return Kill(f)
}

func KillByNameAndPort(name string, portNum int) (int, int, error) {
	f, err := newFilter(name, portNum, true)
	if err != nil {
		return 0, 0, err
	}
	return Kill(f)
}
