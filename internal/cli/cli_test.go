package cli

import (
	"bytes"
	"strings"
	"testing"
)

func execute(t *testing.T, args ...string) (string, error) {
	t.Helper()
	out := &bytes.Buffer{}
	root := NewRootCmd()
	root.SetOut(out)
	root.SetErr(&bytes.Buffer{})
	root.SetArgs(args)
	err := root.Execute()
	return out.String(), err
}

func TestExecute(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		wantErr     bool
		errContains string
		wantOut     string
	}{
		{"версия", []string{"--version"}, false, "", "version"},
		{"help", []string{"--help"}, false, "", "Usage:"},
		{"неизвестная команда", []string{"foo"}, true, "unknown command", ""},
		{"kill без флагов", []string{"kill"}, true, "kill requires --name or --port", ""},
		{"list -port нечисловой", []string{"list", "--port", "abc"}, true, "invalid argument", ""},
		{"kill -port нечисловой", []string{"kill", "--port", "abc"}, true, "invalid argument", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := execute(t, tt.args...)

			if (err != nil) != tt.wantErr {
				t.Fatalf("ошибка = %v, ожидалось wantErr = %v", err, tt.wantErr)
			}
			if tt.wantErr && !strings.Contains(err.Error(), tt.errContains) {
				t.Errorf("ошибка = %q, должна содержать %q", err, tt.errContains)
			}
			if tt.wantOut != "" && !strings.Contains(out, tt.wantOut) {
				t.Errorf("вывод %q должен содержать %q", out, tt.wantOut)
			}
		})
	}
}

func TestExecuteListNoFlags(t *testing.T) {
	out, err := execute(t, "list")
	if err != nil {
		t.Fatalf("ошибка = %v, ожидался nil", err)
	}
	for _, want := range []string{"PROCESS", "PORT", "PID"} {
		if !strings.Contains(out, want) {
			t.Errorf("вывод %q должен содержать заголовок таблицы %q", out, want)
		}
	}
}
