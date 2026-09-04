package main

import (
	"fmt"
	"os"

	"github.com/Awak3r/PortKiller/internal/cli"
)

func main() {
	err := cli.Run(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
