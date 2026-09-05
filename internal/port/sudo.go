package port

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

// ErrNeedRoot reports that the current command requires root privileges
// and the process is not running as root. The caller (main) decides how
// to escalate — this keeps os.Exit/exec out of library code, so tests
// never trigger a real sudo restart (review item 8).
var ErrNeedRoot = errors.New("root privileges required")

// rootArgs returns the argv for the escalated re-execution: the absolute
// path to the current executable plus the original command-line arguments.
func rootArgs() ([]string, error) {
	exe, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("can't locate own executable: %w", err)
	}
	return append([]string{exe}, os.Args[1:]...), nil
}

// RunAsRoot re-executes the utility under sudo and never returns:
//   - child exited normally -> the parent exits with the child's exit code;
//   - child was killed by a signal -> the parent exits with 1;
//   - sudo itself failed (password declined, no sudo, no terminal) -> error.
//
// EnsureRoot below is the thin wrapper used by main; tests should rely on
// RequireRoot/ErrNeedRoot instead of calling this.
func RunAsRoot() error {
	args, err := rootArgs()
	if err != nil {
		return err
	}

	cmd := exec.Command("sudo", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()

	var exitErr *exec.ExitError
	switch {
	case errors.As(err, &exitErr):
		code := exitErr.ExitCode()
		if code < 0 {
			// процесс убит сигналом — осмысленного кода выхода нет
			code = 1
		}
		os.Exit(code)
		return nil // unreachable: os.Exit does not return
	case err != nil:
		return fmt.Errorf("sudo: %w", err)
	default:
		os.Exit(0)
		return nil // unreachable
	}
}

// RequireRoot reports ErrNeedRoot if the process lacks root privileges.
// It performs no side effects and is safe to call from library code.
func RequireRoot() error {
	if os.Geteuid() == 0 {
		return nil
	}
	return ErrNeedRoot
}

// EnsureRoot escalates when needed; under a non-root process it never
// returns normally (the process is replaced by the sudo run or exits).
// Kept for main: it preserves the previous one-call ergonomics while
// isolating os.Exit inside a single, documented place.
func EnsureRoot() error {
	if os.Geteuid() == 0 {
		return nil
	}
	if err := RunAsRoot(); err != nil {
		fmt.Fprintln(os.Stderr, "portkiller:", err)
		os.Exit(1)
	}
	return nil // unreachable
}
