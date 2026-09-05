package cli

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Awak3r/PortKiller/internal/commands"
	"github.com/Awak3r/PortKiller/internal/port"
)

var errKillRequiresFilter = errors.New("kill requires -name or -port")

// escalate re-runs the utility under sudo if the current process is not root.
// Called only from kill commands, after flag validation (review item 5):
// an invalid port must fail before the user is asked for a password.
func escalate() error {
	return port.EnsureRoot()
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

			// валидация до эскалации: невалидный порт не должен спрашивать пароль
			if portSet && (port < 1 || port > 65535) {
				return fmt.Errorf("invalid port (1-65535)")
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

// doKill runs the kill function and reports found/killed stats
// even when some kills failed (partial success).
func doKill(kill func() (int, int, error)) error {
	found, killed, err := kill()
	if found > 0 {
		fmt.Printf("found %d process(es), killed %d\n", found, killed)
	}
	return err
}
