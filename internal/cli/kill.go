package cli

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Awak3r/PortKiller/internal/commands"
)

var errKillRequiresFilter = errors.New("kill requires -name or -port")

func newKillCmd() *cobra.Command {
	var port int
	var name string

	cmd := &cobra.Command{
		Use:   "kill",
		Short: "Kill processes by name and/or port",

		RunE: func(cmd *cobra.Command, args []string) error {
			nameSet := cmd.Flags().Changed("name")
			portSet := cmd.Flags().Changed("port")

			switch {
			case nameSet && portSet:
				return doKill(func() (int, int, error) {
					return commands.KillByNameAndPort(name, port)
				})
			case nameSet:
				return doKill(func() (int, int, error) {
					return commands.KillByName(name)
				})
			case portSet:
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

// doKill runs the kill function and reports found/killed stats
// even when some kills failed (partial success).
func doKill(kill func() (int, int, error)) error {
	found, killed, err := kill()
	if found > 0 {
		fmt.Printf("found %d process(es), killed %d\n", found, killed)
	}
	return err
}
