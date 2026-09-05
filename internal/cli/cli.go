package cli

// Execute runs the root command with the given arguments.
func Execute(args []string) error {
	root := NewRootCmd()
	root.SetArgs(args)
	return root.Execute()
}
