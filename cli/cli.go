package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/Awak3r/PortKiller/commands"
	"github.com/Awak3r/PortKiller/internal/version"
	"github.com/Awak3r/PortKiller/process"
)

func ArgParse() error {
	if len(os.Args) < 2 {
		fmt.Println("run with --help for usage")
		return nil
	}
	switch os.Args[1] {
	case "--version", "-v", "version":
		fmt.Println(version.Full())

	case "--help", "-help", "-h":
		commands.PrintUsage()

	case "list":
		proc, err := process.Collect()
		if err != nil {
			return fmt.Errorf("error listing processes: %w", err)
		}
		cmd := flag.NewFlagSet("list", flag.ExitOnError)
		name := cmd.String("name", "", "process name")
		port := cmd.Int("port", 0, "port")
		cmd.Parse(os.Args[2:])

		switch {
		case *name != "":
			commands.ListByName(*name, proc)
		case *port != 0:
			commands.ListByPort(*port, proc)
		default:
			commands.ListAll(proc)
		}

	case "kill":
		proc, err := process.Collect()
		if err != nil {
			return fmt.Errorf("error listing processes: %w", err)
		}
		cmd := flag.NewFlagSet("kill", flag.ExitOnError)
		name := cmd.String("name", "", "process name")
		port := cmd.Int("port", 0, "port")
		cmd.Parse(os.Args[2:])
		switch {
		case *name != "":
			found, killed, err := commands.KillByName(*name, proc)
			fmt.Printf("Found %d processes matching '%s'. Successfully terminated: %d\n", *name, found, killed)
			if err != nil {
				return fmt.Errorf("error terminating processes matching '%s': %w", *name, err)
			}

		case *port != 0:
			found, killed, err := commands.KillByPort(*port, proc)
			fmt.Printf("Found %d processes on port %d. Successfully terminated: %d\n", *port, found, killed)
			if err != nil {
				return fmt.Errorf("error terminating processes on port %d: %w", *port, err)
			}
		default:
			fmt.Println("specify -name or -port")
		}
	default:
		fmt.Println("unknown command:", os.Args[1])
		commands.PrintUsage()
	}
	return nil
}
