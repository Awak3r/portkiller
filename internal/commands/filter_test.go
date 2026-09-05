package commands

import "testing"

func TestProcFilterMatch(t *testing.T) {
	tests := []struct {
		name     string
		filter   procFilter
		procName string
		procPort int
		want     bool
	}{
		{"пустой фильтр матчит всё", procFilter{}, "nginx", 80, true},
		{"имя подстрокой", procFilter{name: "ngin"}, "nginx", 80, true},
		{"имя регистронезависимо", procFilter{name: "NGIN"}, "nginx", 80, true},
		{"имя не совпало", procFilter{name: "node"}, "nginx", 80, false},
		{"порт совпал", procFilter{port: 80}, "nginx", 80, true},
		{"порт не совпал", procFilter{port: 443}, "nginx", 80, false},
		{"имя+порт ок", procFilter{name: "nginx", port: 80}, "nginx", 80, true},
		{"имя+порт порт мимо", procFilter{name: "nginx", port: 443}, "nginx", 80, false},
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
	if _, err := newFilter("", 70000, true); err == nil {
		t.Error("порт 70000 должен давать ошибку")
	}
	if _, err := newFilter("", 0, true); err == nil {
		t.Error("порт 0 при явном -port должен давать ошибку")
	}
	if _, err := newFilter("nginx", 8080, true); err != nil {
		t.Errorf("валидный порт не должен давать ошибку: %v", err)
	}
	// имя-фильтр без порта — валидация не нужна
	if _, err := newFilter("nginx", 0, false); err != nil {
		t.Errorf("фильтр по имени не должен давать ошибку: %v", err)
	}
}
