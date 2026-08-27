package cli

import (
	"fmt"
	"os"
	"strconv"

	"github.com/shirou/gopsutil/v4/process"
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
		return nil

	default:
		if len(os.Args) < 3 {
			return fmt.Errorf("unknown argument %q, use --help for usage", os.Args[1])
		}
		switch os.Args[2] {
		case "list":
			p, err := process.Collect()
			if err != nil {
				return err
			}
			name := ""
			port := 0
			for i := 3; i < len(os.Args); i++ {
				switch os.Args[i] {
				case "-name":
					i++
					name = os.Args[i]
				case "-port":
					i++
					port, _ = strconv.Atoi(os.Args[i])
				}
			}
			if name != "" {
				commands.ListByName(name, p)
			} else if port != 0 {
				commands.ListByPort(port, p)
			} else {
				commands.ListAll(p)
			}
			return nil

		case "kill":
			p, err := process.Collect()
			if err != nil {
				return err
			}
			name := ""
			port := 0
			for i := 3; i < len(os.Args); i++ {
				switch os.Args[i] {
				case "-name":
					i++
					name = os.Args[i]
				case "-port":
					i++
					port, _ = strconv.Atoi(os.Args[i])
				}
			}
			if name != "" {
				found, killed, err := commands.KillByName(name, p)
				if err != nil {
					return err
				}
				fmt.Printf("found %d process(es), killed %d\n", found, killed)
			} else if port != 0 {
				found, killed, err := commands.KillByPort(port, p)
				if err != nil {
					return err
				}
				fmt.Printf("found %d process(es), killed %d\n", found, killed)
			} else {
				return fmt.Errorf("kill requires -name or -port")
			}
			return nil

		default:
			return fmt.Errorf("unknown command %q, use --help for usage", os.Args[2])
		}
	}
}
