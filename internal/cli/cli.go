package cli

func Execute(args []string) error {
	root := NewRootCmd()
	root.SetArgs(args)
	return root.Execute()
}
