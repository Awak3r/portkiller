package version

import "fmt"

var (
	Version = "dev"
	Commit  = "unknown"
	Date    = "unknown"
)

func Full() string {
	return fmt.Sprintf(
		"%s (commit %s, built %s)",
		Version,
		Commit,
		Date,
	)
}
