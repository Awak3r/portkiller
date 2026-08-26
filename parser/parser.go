package parser

import (
	"flag"
	"fmt"
	"os"

	"github.com/Awak3r/PortKiller/commands"
	"github.com/Awak3r/PortKiller/internal/version"
	"github.com/Awak3r/PortKiller/utils"
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
        proc := utils.Collect()
		cmd := flag.NewFlagSet("list", flag.ExitOnError)
		name := cmd.String("name", "", "имя процесса")
		port := cmd.Int("port", 0, "порт")
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
        proc := utils.Collect()
		cmd := flag.NewFlagSet("kill", flag.ExitOnError)
		name := cmd.String("name", "", "имя процесса")
		port := cmd.Int("port", 0, "порт")
		cmd.Parse(os.Args[2:])
		switch {
            case *name != "":
                found, killed, err := commands.KillByName(*name, proc)
                fmt.Printf("Найдено процессов '%s': %d. Успешно завершено: %d\n", *name, found, killed)
                if err != nil {
                    return fmt.Errorf("ошибка при завершении процессов '%s': %w", *name, err)
                }

            case *port != 0:
                found, killed, err := commands.KillByPort(*port, proc)
                fmt.Printf("Найдено процессов на порту %d: %d. Успешно завершено: %d\n", *port, found, killed)
                if err != nil {
                    return fmt.Errorf("ошибка при завершении процессов на порту %d: %w", *port, err)
            }
		default:
			fmt.Println("укажи -name или -port")
		}
	default:
		fmt.Println("неизвестная команда:", os.Args[1])
		commands.PrintUsage()
	}
    return nil
}
