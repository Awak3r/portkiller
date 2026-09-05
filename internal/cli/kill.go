package cli

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Awak3r/PortKiller/internal/commands"
	"github.com/Awak3r/PortKiller/internal/port"
)

var errKillRequiresFilter = errors.New("kill requires --name or --port")

func newKillCmd() *cobra.Command {
	var portFlag int
	var name string
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "kill",
		Short: "Kill processes by name and/or port",

		RunE: func(cmd *cobra.Command, args []string) error {
			nameSet := cmd.Flags().Changed("name")
			portSet := cmd.Flags().Changed("port")

			if portSet {
				if err := commands.ValidatePort(portFlag); err != nil {
					return err
				}
			}

			var portPtr *int
			if portSet {
				portPtr = &portFlag
			}
			filter, err := commands.NewFilter(name, portPtr)
			if err != nil {
				return err
			}

			if !nameSet && !portSet {
				return errKillRequiresFilter
			}

			if dryRun {
				return filter.List(cmd.OutOrStdout())
			}

			if err := port.RequireRoot(); err != nil {
				return err
			}

			return doKill(cmd, func() (int, int, error) {
				return filter.Kill(cmd.Context())
			})
		},
	}

	cmd.Flags().IntVarP(&portFlag, "port", "p", 0, "port to kill by")
	cmd.Flags().StringVarP(&name, "name", "n", "", "process name to kill by")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "show what would be killed without killing")

	return cmd
}

func doKill(cmd *cobra.Command, kill func() (int, int, error)) error {
	found, killed, err := kill()
	if found > 0 {
		fmt.Fprintf(cmd.OutOrStdout(), "found %d process(es), killed %d\n", found, killed)
	}
	return err
}
