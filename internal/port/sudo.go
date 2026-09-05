package port

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

// EnsureRoot перезапускает утилиту через sudo, если текущий процесс не под root.
//
// Родительский процесс после запуска ребёнка всегда завершается:
//   - ребёнок отработал успешно -> exit 0;
//   - ребёнок завершился сам -> пробрасывается его exit code;
//   - sudo не смог запустить ребёнка (отменён пароль, нет sudo) -> stderr + exit 1.
//
// Управление в родительский процесс не возвращается никогда,
// поэтому двойное выполнение команды (баг из ревью №1) невозможно.
func EnsureRoot() {
	if os.Geteuid() == 0 {
		return
	}

	exe, err := os.Executable()
	if err != nil {
		fmt.Fprintln(os.Stderr, "portkiller: can't locate own executable:", err)
		os.Exit(1)
	}

	cmd := exec.Command("sudo", append([]string{exe}, os.Args[1:]...)...)
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
	case err != nil:
		fmt.Fprintln(os.Stderr, "portkiller: sudo failed:", err)
		os.Exit(1)
	default:
		os.Exit(0)
	}
}
