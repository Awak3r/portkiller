package commands

import "testing"

func TestNameMatches(t *testing.T) {
	tests := []struct {
		proc    string
		pattern string
		want    bool
	}{
		{"nginx", "ngin", true},
		{"NGINX", "ngin", true},
		{"nginx-helper", "nginx", true},
		{"nginx", "NGINX", true},
		{"node", "nginx", false},
		{"", "ngin", false},
		{"nginx", "", true}, // пустой шаблон матчит всё — отсекается на уровне флагов
	}

	for _, tt := range tests {
		if got := NameMatches(tt.proc, tt.pattern); got != tt.want {
			t.Errorf("NameMatches(%q, %q) = %v, want %v", tt.proc, tt.pattern, got, tt.want)
		}
	}
}
