[English](README.md) | [Русский](README.ru.md)

# PortKiller 🔪

A small CLI utility to see which processes occupy TCP ports and to terminate them — by port, by name, or by both. Single static binary, no runtime dependencies.

Built with [cobra](https://github.com/spf13/cobra) and [gopsutil](https://github.com/shirou/gopsutil).

## Features

- 📋 `list` — table of processes listening on TCP ports (works **without** root)
- 💀 `kill` — terminate processes by port, by name, or by name + port (asks for root)
- 🔍 Name matching: case-insensitive **substring** — the same semantics for `list` and `kill`, so what you see in `list` is exactly what `kill` will terminate
- 🪜 Two-stage termination: SIGTERM first, SIGKILL fallback for processes that ignore it (see [How killing works](#how-killing-works))
- 📊 Kill report even on partial failure: `found N process(es), killed M`
- ⚡ Root escalation only when needed: `list` never prompts for a password; `kill` validates flags **before** asking for one

## Installation

### Method 1: via Go

Requires **Go 1.26+** (see `go.mod`):

```bash
go install github.com/Awak3r/PortKiller/cmd/portkiller@latest
```

Make sure `$(go env GOPATH)/bin` is in your `PATH` (one-time setup):

```bash
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.zshrc   # or ~/.bashrc
source ~/.zshrc
```

### Method 2: pre-built binaries

Grab an archive from [Releases](https://github.com/Awak3r/PortKiller/releases) — currently Linux `amd64` / `arm64` (`portkiller-linux-*`), see the release assets for what is available. Extract and move to a directory in your `PATH`:

```bash
tar -xzf portkiller-linux-amd64.tar.gz
sudo mv portkiller /usr/local/bin/
```

> Releases are built by GoReleaser on tag push (`v*`) with version/commit/date embedded via `-ldflags`.

### Method 3: build from source

```bash
git clone https://github.com/Awak3r/PortKiller.git
cd PortKiller
go build -o portkiller ./cmd/portkiller
```

## Usage

```text
portkiller — kill processes by name or port

Usage:
  portkiller [command]

Available Commands:
  list   List processes
  kill   Kill processes by name and/or port

Flags:
  -h, --help      help for portkiller
  -v, --version   version for portkiller
```

The binary name is `portkiller`. Commands use GNU-style long flags (`--port`, `--name`) with single-letter aliases (`-p`, `-n`); both single-dash long forms and double-dash forms are accepted by the flag parser.

### List occupied ports (no root needed)

```bash
portkiller list
```

```text
PROCESS     PORT   PID
-------     ----   ----
node        3000   12345
postgres    5432   999
```

Filter the table the same way as `kill`:

```bash
portkiller list --name node        # by name (substring, case-insensitive)
portkiller list --port 3000        # by port
portkiller list --name node --port 3000
```

### Kill by port

```bash
sudo portkiller kill --port 5000
```

### Kill by name

```bash
sudo portkiller kill --name node
```

Matches every process whose name contains `node` (case-insensitive) — `node`, `NODE`, `my-node-app`.

### Kill by name and port

```bash
sudo portkiller kill --name node --port 3000
```

Both conditions must match at once.

### Exit codes

| Code | Meaning |
|------|---------|
| `0`  | Success |
| `1`  | Invalid usage, nothing found, or the process(es) could not be killed |

## Flags

| Flag (long / short) | Applies to | Description |
|---------------------|------------|-------------|
| `--port` / `-p`     | `list`, `kill` | TCP port, integer 1–65535 |
| `--name` / `-n`     | `list`, `kill` | Process name, case-insensitive substring |
| `--help` / `-h`     | everywhere | Help |
| `--version` / `-v`  | root command | Version: `portkiller version X (commit Y, built Z)` |

Validation: an invalid port (`portkiller kill --port 70000`) fails **immediately, before any sudo prompt**; `kill` without flags prints `kill requires -name or -port`.

## Root & permissions

- **`list` does not require root.** Reading the TCP table is available to any user (you may see fewer details for processes owned by other users, but ports and PIDs are shown).
- **`kill` requires root** to send signals to processes you don't own. When you run `portkiller kill ...` as a regular user:
  1. your flags are validated first (invalid port → error, no password asked);
  2. the tool re-executes itself under `sudo` and asks for your password;
  3. the child's exit code is propagated back, so scripts can rely on it.
- Run it directly under root/sudo if you prefer: `sudo portkiller kill --name nginx`.

## How killing works

The termination is a **two-stage escalation** (implemented in `internal/port.KillByPid`):

**Stage 1 — SIGTERM.** The process receives a "polite" termination request. A well-behaved program may intercept it and shut down gracefully: close sockets, flush buffers, remove pidfiles, save state.

**Grace period.** PortKiller then waits up to **2 seconds** (polling every 100 ms) for the process to actually exit.

**Stage 2 — SIGKILL fallback.** If the process is still alive after the grace period (hung, stuck in I/O, or explicitly ignoring/trapping SIGTERM), PortKiller sends **SIGKILL**. This signal cannot be caught, blocked, or ignored — the kernel destroys the process immediately, without any cleanup. That is the fallback that guarantees a stuck process still dies.

Why bother with stage 1 at all? Because SIGKILL takes away the process's chance to clean up after itself — most of the time SIGTERM is enough and is the safe choice; SIGKILL is the last resort. Concretely:

| Scenario | What happens |
|----------|--------------|
| Normal app (nginx, node, postgres) | Dies on SIGTERM, usually in milliseconds |
| App that traps SIGTERM and ignores it | Gets SIGKILL after the 2 s grace period |
| Process that exited between listing and the signal (`ESRCH`) | Not an error — the goal is already achieved |
| `SIGKILL` itself failed (e.g. PID reuse edge case) | Returned as an error, reported in the kill summary |

The per-process result is aggregated: even if some kills fail, you still get `found N process(es), killed M` plus the joined errors, and the exit code is non-zero.

## Cross-platform status

`list`/`kill` logic is cross-platform via gopsutil, but **today the release pipeline builds Linux binaries only** (`amd64`, `arm64`), and root escalation uses `sudo`, which is a POSIX convention (on Windows it would need a different mechanism). See `.goreleaser.yaml` — adding `darwin`/`windows` targets is on the roadmap.

## Development

```bash
go build ./...        # build
go vet ./...          # static analysis
go test ./...         # unit tests
gofmt -l .            # formatting check
```

Project layout:

```text
cmd/portkiller/       main: os.Exit + sudo escalation (the only place allowed to exit)
internal/cli/         cobra commands and flag parsing
internal/commands/    list/kill pipeline: Collect -> filter -> action
internal/port/        gopsutil collection, two-stage kill, root helpers
internal/version/     build metadata (injected via -ldflags)
```

## License

Apache-2.0 — see [LICENSE](LICENSE).
