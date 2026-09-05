package commands

import (
	"context"
	"fmt"
	"io"
)

func PrintUsage() {
	fmt.Println(`portkiller — kill processes by name or port

Usage:
  portkiller list [-name NAME | -port PORT]
  portkiller kill [-name NAME | -port PORT]
  portkiller -version
  portkiller -help`)
}

func ListAll(w io.Writer) error {
	return List(w, "", 0, false)
}

func ListByName(w io.Writer, name string) error {
	return List(w, name, 0, false)
}

func ListByPort(w io.Writer, portNum int) error {
	return List(w, "", portNum, true)
}

func ListByNameAndPort(w io.Writer, name string, portNum int) error {
	return List(w, name, portNum, true)
}

func KillByName(ctx context.Context, name string) (int, int, error) {
	return Kill(ctx, procFilter{name: name})
}

func KillByPort(ctx context.Context, portNum int) (int, int, error) {
	f, err := newFilter("", portNum, true)
	if err != nil {
		return 0, 0, err
	}
	return Kill(ctx, f)
}

func KillByNameAndPort(ctx context.Context, name string, portNum int) (int, int, error) {
	f, err := newFilter(name, portNum, true)
	if err != nil {
		return 0, 0, err
	}
	return Kill(ctx, f)
}
