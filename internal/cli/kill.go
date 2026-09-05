package cli

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Awak3r/PortKiller/internal/commands"
	"github.com/Awak3r/PortKiller/internal/port"
)

var errKillRequiresFilter = errors.New("kill requires -name or -port")

func escalate() error {
	return port.RequireRoot()
}

func newKillCmd() *cobra.Command {
	var port int
	var name string

	cmd := &cobra.Command{
		Use:   "kill",
		Short: "Kill processes by name and/or port",

		RunE: func(cmd *cobra.Command, args []string) error {
			nameSet := cmd.Flags().Changed("name")
			portSet := cmd.Flags().Changed("port")

			if portSet {
				if err := commands.ValidatePort(port); err != nil {
					return err
				}
			}

			switch {
			case nameSet && portSet:
				if err := escalate(); err != nil {
					return err
				}
				return doKill(func() (int, int, error) {
					return commands.KillByNameAndPort(name, port)
				})
			case nameSet:
				if err := escalate(); err != nil {
					return err
				}
				return doKill(func() (int, int, error) {
					return commands.KillByName(name)
				})
			case portSet:
				if err := escalate(); err != nil {
					return err
				}
				return doKill(func() (int, int, error) {
					return commands.KillByPort(port)
				})
			default:
				return errKillRequiresFilter
			}
		},
	}

	cmd.Flags().IntVarP(
		&port,
		"port",
		"p",
		0,
		"port to kill by",
	)

	cmd.Flags().StringVarP(
		&name,
		"name",
		"n",
		"",
		"process name to kill by",
	)

	return cmd
}

func doKill(kill func() (int, int, error)) error {
	found, killed, err := kill()
	if found > 0 {
		fmt.Printf("found %d process(es), killed %d\n", found, killed)
	}
	return err
}
