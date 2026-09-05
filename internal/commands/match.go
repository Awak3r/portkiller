package commands

import "strings"

func NameMatches(procName, pattern string) bool {
	return strings.Contains(
		strings.ToLower(procName),
		strings.ToLower(pattern),
	)
}
