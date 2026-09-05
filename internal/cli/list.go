package cli

import (
	"github.com/spf13/cobra"

	"github.com/Awak3r/PortKiller/internal/commands"
)

func newListCmd() *cobra.Command {
	var portFlag int
	var name string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List processes",

		RunE: func(cmd *cobra.Command, args []string) error {
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
			return filter.List(cmd.OutOrStdout())
		},
	}

	cmd.Flags().IntVarP(&portFlag, "port", "p", 0, "port to list")
	cmd.Flags().StringVarP(&name, "name", "n", "", "name to list")

	return cmd
}
