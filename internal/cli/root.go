package cli

import (
	"github.com/Awak3r/PortKiller/internal/version"
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:           "portkiller",
		Short:         "portkiller — kill processes by name or port",
		Version:       version.Full(),
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	rootCmd.SetVersionTemplate("{{with .Name}}{{printf \"%s \" .}}{{end}}{{printf \"version %s\" .Version}}\n")
	rootCmd.AddCommand(newListCmd(), newKillCmd())
	return rootCmd
}
