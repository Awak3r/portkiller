package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"

	"github.com/Awak3r/PortKiller/internal/commands"
	"github.com/Awak3r/PortKiller/internal/version"
)

func Run(args []string) error {
	if len(args) < 2 {
		fmt.Println("run with --help for usage")
		return nil
	}
	switch args[1] {
	case "--version", "-v", "version":
		fmt.Println(version.Full())
		return nil

	case "--help", "-help", "-h":
		commands.PrintUsage()
		return nil

	case "list":
		return runList(args[2:])

	case "kill":
		return runKill(args[2:])

	default:
		return fmt.Errorf("unknown command %q, use --help for usage", args[1])
	}
}

func parseFlags(fs *flag.FlagSet, args []string) (map[string]bool, int, string, error) {
	fs.SetOutput(io.Discard)
	port := fs.Int("port", 0, "port to list")
	name := fs.String("name", "", "name to list")
	if err := fs.Parse(args); err != nil {
		return nil, 0, "", err
	}
	var flagsSet = make(map[string]bool, 2)
	fs.Visit(func(f *flag.Flag) {
		flagsSet[f.Name] = true
	})
	return flagsSet, *port, *name, nil
}

func runList(args []string) error {
	fs := flag.NewFlagSet("list", flag.ContinueOnError)
	flagsSet, port, name, err := parseFlags(fs, args)
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			commands.PrintUsage()
			return nil
		}
		return fmt.Errorf("list: %w", err)
	}
	if flagsSet["name"] && flagsSet["port"] {
		return commands.ListByNameAndPort(name, port)
	} else if flagsSet["name"] {
		return commands.ListByName(name)
	} else if flagsSet["port"] {
		return commands.ListByPort(port)
	} else {
		commands.ListAll()
	}
	return nil
}

func runKill(args []string) error {
	fs := flag.NewFlagSet("kill", flag.ContinueOnError)
	flagsSet, port, name, err := parseFlags(fs, args)
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			commands.PrintUsage()
			return nil
		}
		return fmt.Errorf("kill: %w", err)
	}
	if flagsSet["name"] && flagsSet["port"] {
		found, killed, err := commands.KillByNameAndPort(name, port)
		if err != nil {
			return err
		}
		fmt.Printf("found %d process(es), killed %d\n", found, killed)
	} else if flagsSet["name"] {
		found, killed, err := commands.KillByName(name)
		if err != nil {
			return err
		}
		fmt.Printf("found %d process(es), killed %d\n", found, killed)
	} else if flagsSet["port"] {
		found, killed, err := commands.KillByPort(port)
		if err != nil {
			return err
		}
		fmt.Printf("found %d process(es), killed %d\n", found, killed)
	} else {
		return fmt.Errorf("kill requires -name or -port")
	}
	return nil
}
