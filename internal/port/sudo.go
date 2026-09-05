package port

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

var ErrNeedRoot = errors.New("root privileges required")

func rootArgs() ([]string, error) {
	exe, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("can't locate own executable: %w", err)
	}
	return append([]string{exe}, os.Args[1:]...), nil
}

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
			code = 1
		}
		os.Exit(code)
		return nil
	case err != nil:
		return fmt.Errorf("sudo: %w", err)
	default:
		os.Exit(0)
		return nil
	}
}

func RequireRoot() error {
	if os.Geteuid() == 0 {
		return nil
	}
	return ErrNeedRoot
}

// HasTTY reports whether the process has a terminal attached, which
// sudo needs to prompt for a password.
func HasTTY() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}

func EnsureRoot() error {
	if os.Geteuid() == 0 {
		return nil
	}
	if err := RunAsRoot(); err != nil {
		fmt.Fprintln(os.Stderr, "portkiller:", err)
		os.Exit(1)
	}
	return nil
}
