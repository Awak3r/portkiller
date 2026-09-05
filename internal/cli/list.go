package cli

import (
	"github.com/spf13/cobra"

	"github.com/Awak3r/PortKiller/internal/commands"
)

func newListCmd() *cobra.Command {
	var port int
	var name string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List processes",

		RunE: func(cmd *cobra.Command, args []string) error {
			nameSet := cmd.Flags().Changed("name")
			portSet := cmd.Flags().Changed("port")

			// io.Writer от cobra: таблица попадает в тот же поток,
			// что и весь вывод команды (и в тестовые буферы)
			out := cmd.OutOrStdout()

			if nameSet && portSet {
				return commands.ListByNameAndPort(out, name, port)
			}
			if nameSet {
				return commands.ListByName(out, name)
			}
			if portSet {
				return commands.ListByPort(out, port)
			}
			return commands.ListAll(out)
		},
	}

	cmd.Flags().IntVarP(
		&port,
		"port",
		"p",
		0,
		"port to list",
	)

	cmd.Flags().StringVarP(
		&name,
		"name",
		"n",
		"",
		"name to list",
	)

	return cmd
}
