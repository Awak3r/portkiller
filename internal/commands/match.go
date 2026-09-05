package commands

import "strings"

// NameMatches reports whether a process name matches the user pattern.
//
// Matching is a case-insensitive substring ("ngin" matches "nginx",
// "NGINX-Helper", "innginx"). This is the single matching semantics shared
// by `list` and `kill` (CODE_REVIEW.md, item 4): what `list` shows is
// exactly what `kill` will terminate. An empty pattern matches everything;
// the flag layer rejects empty names before reaching this function.
func NameMatches(procName, pattern string) bool {
	return strings.Contains(
		strings.ToLower(procName),
		strings.ToLower(pattern),
	)
}
