package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/Awak3r/PortKiller/internal/cli"
	"github.com/Awak3r/PortKiller/internal/port"
)

func main() {
	err := cli.Execute(os.Args[1:])
	if err == nil {
		return
	}
	if errors.Is(err, port.ErrNeedRoot) {
		if escErr := port.EnsureRoot(); escErr != nil {
			fmt.Fprintln(os.Stderr, escErr)
			os.Exit(1)
		}
		return
	}
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
