package commands

import "testing"

func TestFilterMatch(t *testing.T) {
	port80 := 80
	port443 := 443
	tests := []struct {
		name     string
		filter   Filter
		procName string
		procPort int
		want     bool
	}{
		{"empty filter matches all", Filter{}, "nginx", 80, true},
		{"name substring", Filter{Name: "ngin"}, "nginx", 80, true},
		{"name case-insensitive", Filter{Name: "NGIN"}, "nginx", 80, true},
		{"name mismatch", Filter{Name: "node"}, "nginx", 80, false},
		{"port match", Filter{Port: &port80}, "nginx", 80, true},
		{"port mismatch", Filter{Port: &port443}, "nginx", 80, false},
		{"name+port ok", Filter{Name: "nginx", Port: &port80}, "nginx", 80, true},
		{"name+port port off", Filter{Name: "nginx", Port: &port443}, "nginx", 80, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.filter.match(ProcessInfo{Name: tt.procName, Port: tt.procPort, Pid: 1})
			if got != tt.want {
				t.Errorf("match(%q, %d) = %v, want %v", tt.procName, tt.procPort, got, tt.want)
			}
		})
	}
}

func TestNewFilterValidation(t *testing.T) {
	bad := 70000
	zero := 0
	good := 8080

	if _, err := NewFilter("", &bad); err == nil {
		t.Error("port 70000 must be rejected")
	}
	if _, err := NewFilter("", &zero); err == nil {
		t.Error("explicit port 0 must be rejected")
	}
	if _, err := NewFilter("nginx", &good); err != nil {
		t.Errorf("valid port must not fail: %v", err)
	}
	if _, err := NewFilter("nginx", nil); err != nil {
		t.Errorf("name-only filter must not fail: %v", err)
	}
}
