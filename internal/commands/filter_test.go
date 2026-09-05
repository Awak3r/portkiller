package commands

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/Awak3r/PortKiller/internal/port"
)

type fakeCollector struct {
	procs []port.ProcessInfo
	err   error
}

func (f *fakeCollector) Collect() ([]port.ProcessInfo, error) {
	return f.procs, f.err
}

type fakeKiller struct {
	fail map[int]error
}

func (f *fakeKiller) KillByPid(pid int32) error {
	if err, ok := f.fail[int(pid)]; ok {
		return err
	}
	return nil
}

func TestFilterList(t *testing.T) {
	collector := &fakeCollector{procs: []port.ProcessInfo{
		{Name: "nginx", Pid: 1, Port: 80},
		{Name: "node", Pid: 2, Port: 3000},
	}}
	port3000 := 3000
	f, err := NewFilter("", &port3000, WithCollector(collector))
	if err != nil {
		t.Fatalf("filter: %v", err)
	}

	var out strings.Builder
	if err := f.List(&out); err != nil {
		t.Fatalf("List: %v", err)
	}
	if !strings.Contains(out.String(), "node") || !strings.Contains(out.String(), "3000") {
		t.Errorf("output %q must contain node/3000", out.String())
	}
	if strings.Contains(out.String(), "nginx") {
		t.Errorf("output %q must not contain nginx", out.String())
	}
}

func TestFilterListEmptyNoError(t *testing.T) {
	collector := &fakeCollector{procs: []port.ProcessInfo{
		{Name: "nginx", Pid: 1, Port: 80},
	}}
	f, _ := NewFilter("zzz", nil, WithCollector(collector))

	var out strings.Builder
	if err := f.List(&out); err != nil {
		t.Fatalf("err = %v, want nil (empty result is not an error for list)", err)
	}
	if !strings.Contains(out.String(), "PROCESS") {
		t.Errorf("output %q must contain the table header", out.String())
	}
}

func TestFilterListCollectError(t *testing.T) {
	collector := &fakeCollector{err: errors.New("boom")}
	f, _ := NewFilter("", nil, WithCollector(collector))

	var out strings.Builder
	if err := f.List(&out); err == nil || !strings.Contains(err.Error(), "boom") {
		t.Fatalf("err = %v, want boom", err)
	}
}

func TestFilterKillNoMatch(t *testing.T) {
	collector := &fakeCollector{procs: []port.ProcessInfo{
		{Name: "nginx", Pid: 1, Port: 80},
	}}
	f, _ := NewFilter("zzz", nil, WithCollector(collector), WithKiller(&fakeKiller{}))

	_, _, err := f.Kill(context.Background())
	if !errors.Is(err, ErrProcessNotFound) {
		t.Fatalf("err = %v, want ErrProcessNotFound", err)
	}
}

func TestFilterKillPartialFailure(t *testing.T) {
	collector := &fakeCollector{procs: []port.ProcessInfo{
		{Name: "a", Pid: 1, Port: 100},
		{Name: "b", Pid: 2, Port: 200},
		{Name: "c", Pid: 3, Port: 300},
	}}
	killer := &fakeKiller{fail: map[int]error{2: errors.New("no permission")}}
	f, _ := NewFilter("", nil, WithCollector(collector), WithKiller(killer))

	found, killed, err := f.Kill(context.Background())
	if found != 3 {
		t.Errorf("found = %d, want 3", found)
	}
	if killed != 2 {
		t.Errorf("killed = %d, want 2 (stats must survive partial failure)", killed)
	}
	if err == nil || !strings.Contains(err.Error(), "no permission") {
		t.Errorf("err = %v, want joined no-permission error", err)
	}
}

func TestFilterKillAllSuccess(t *testing.T) {
	collector := &fakeCollector{procs: []port.ProcessInfo{
		{Name: "a", Pid: 1, Port: 100},
		{Name: "b", Pid: 2, Port: 200},
	}}
	f, _ := NewFilter("", nil, WithCollector(collector), WithKiller(&fakeKiller{}))

	found, killed, err := f.Kill(context.Background())
	if err != nil {
		t.Fatalf("err = %v, want nil", err)
	}
	if found != 2 || killed != 2 {
		t.Errorf("found/killed = %d/%d, want 2/2", found, killed)
	}
}

func TestFilterKillForeignProcess(t *testing.T) {
	collector := &fakeCollector{procs: []port.ProcessInfo{
		{Name: "-", Pid: 0, Port: 53},
		{Name: "a", Pid: 1, Port: 100},
	}}
	f, _ := NewFilter("", nil, WithCollector(collector), WithKiller(&fakeKiller{}))

	found, killed, err := f.Kill(context.Background())
	if found != 1 || killed != 1 {
		t.Errorf("found/killed = %d/%d, want 1/1 (pid 0 is skipped, not signaled)", found, killed)
	}
	if err == nil || !strings.Contains(err.Error(), "owned by another user") {
		t.Errorf("err = %v, want foreign-process error", err)
	}
}

func TestFilterKillContextCancelled(t *testing.T) {
	collector := &fakeCollector{procs: []port.ProcessInfo{
		{Name: "a", Pid: 1, Port: 100},
		{Name: "b", Pid: 2, Port: 200},
		{Name: "c", Pid: 3, Port: 300},
	}}
	f, _ := NewFilter("", nil, WithCollector(collector), WithKiller(&fakeKiller{}))

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	found, _, err := f.Kill(ctx)
	if !errors.Is(err, context.Canceled) {
		t.Errorf("err = %v, want context.Canceled", err)
	}
	if found != 0 {
		t.Errorf("found = %d, want 0 (nothing dispatched after cancel)", found)
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
