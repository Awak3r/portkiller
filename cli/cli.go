package cli

import (
	"flag"
	"fmt"
	"github.com/Awak3r/PortKiller/commands"
	"github.com/Awak3r/PortKiller/internal/version"
	"os"
)

func ArgParse() error {
	if len(os.Args) < 2 {
		fmt.Println("run with --help for usage")
		return nil
	}
	switch os.Args[1] {
	case "--version", "-v", "version":
		fmt.Println(version.Full())
		return nil

	case "--help", "-help", "-h":
		commands.PrintUsage()
		return nil

	case "list":
		return runList(os.Args[2:])

	case "kill":
		return runKill(os.Args[2:])

	default:
		return fmt.Errorf("unknown command %q, use --help for usage", os.Args[1])
	}
}

func parseFlags(fs *flag.FlagSet, args []string) (map[string]bool, int, string) {
	port := fs.Int("port", 0, "port to list")
	name := fs.String("name", "", "name to list")
	fs.Parse(args)
	var flagsSet = make(map[string]bool, 2)
	fs.Visit(func(f *flag.Flag) {
		flagsSet[f.Name] = true
	})
	return flagsSet, *port, *name
}

func runList(args []string) error {
	fs := flag.NewFlagSet("list", flag.ExitOnError)
	flagsSet, port, name := parseFlags(fs, args)
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
	fs := flag.NewFlagSet("kill", flag.ExitOnError)
	flagsSet, port, name := parseFlags(fs, args)
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
