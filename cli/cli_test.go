package cli

import (
	"flag"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		wantErr     bool
		errContains string
		blocked     bool
	}{
		{"нет аргументов", []string{"app"}, false, "", false},
		{"версия", []string{"app", "-v"}, false, "", false},
		{"help", []string{"app", "--help"}, false, "", false},
		{"неизвестная команда", []string{"app", "foo"}, true, "unknown command", false},
		{"list без флагов", []string{"app", "list"}, false, "", true},
		{"kill без флагов", []string{"app", "kill"}, true, "requires -name or -port", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.blocked {
				t.Skip("blocked by EnsureRoot/Collect: os.Exit inside library code — unblocks in Step 4")
			}
			err := Run(tt.args)

			if (err != nil) != tt.wantErr {
				t.Fatalf("ошибка = %v, ожидалось wantErr = %v", err, tt.wantErr)
			}
			if tt.wantErr && !strings.Contains(err.Error(), tt.errContains) {
				t.Errorf("ошибка = %q, должна содержать %q", err, tt.errContains)
			}
		})
	}
}
func TestParseFlags(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantPort int
		wantName string
		wantErr  bool
	}{
		{"только порт", []string{"-port", "8080"}, 8080, "", false},
		{"порт и имя", []string{"-name", "nginx", "-port", "80"}, 80, "nginx", false},
		{"ошибка парсинга", []string{"-port", "abc"}, 0, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := flag.NewFlagSet("test", flag.ContinueOnError)
			flagsSet, port, name, err := parseFlags(fs, tt.args)

			if (err != nil) != tt.wantErr {
				t.Fatalf("ошибка = %v, wantErr = %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if port != tt.wantPort || name != tt.wantName {
				t.Errorf("получено port=%d, name=%q; ожидалось port=%d, name=%q",
					port, name, tt.wantPort, tt.wantName)
			}

			if tt.args[0] == "-port" && !flagsSet["port"] {
				t.Error("флаг port должен быть помечен как установленный")
			}
		})
	}
}
