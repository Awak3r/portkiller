[English](README.md) | [Русский](README.ru.md)
# PortKiller 🔪

CLI utility for managing processes via occupied ports. Cross-platform, single binary with zero dependencies.

## Features

- 📋 List all occupied TCP ports with process details
- 💀 Kill a process by port number
- 🔍 Kill processes by name (supports substrings, case-insensitive)

## Installation

### Method 1: Via Go (Recommended)

Requires **Go 1.22+**:

```bash
go install github.com/Awak3r/PortKiller@latest
```

Add Go binaries to PATH (one-time setup):

```bash
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.zshrc   # or ~/.bashrc
source ~/.zshrc
```

### Method 2: Pre-built Binaries (No Go required)

1. Go to [Releases](https://github.com/Awak3r/PortKiller/releases)
2. Download the archive for your OS:
   - Linux: `portkiller-linux-amd64.tar.gz`
   - macOS: `portkiller-darwin-arm64.tar.gz` (Apple Silicon) or `-amd64` (Intel)
   - Windows: `portkiller-windows-amd64.zip`
3. Extract and move to a directory in your PATH:

```bash
# Linux/macOS
sudo mv portkiller /usr/local/bin/
sudo chmod +x /usr/local/bin/portkiller

# Windows
# Copy portkiller.exe to C:\Windows\ or add the folder to PATH
```

### Method 3: Building from Source

```bash
git clone https://github.com/Awak3r/PortKiller.git
cd PortKiller
go build -o portkiller .
sudo mv portkiller /usr/local/bin/
```

## Usage

> ⚠️ **Important:** The utility requires administrator (sudo) privileges — it will automatically prompt for a password upon launch.

### List Occupied Ports

```bash
PortKiller list
```

Outputs a table: `PORT | PID | NAME`

### Kill by Port

```bash
PortKiller kill -port 5000
```

### Kill by Name

```bash
PortKiller kill -name node
```

## Flags

| Flag | Description |
|------|-------------|
| `-port <number>` | Port number to kill (1–65535) |
| `-name <string>` | Process name (substring, case-insensitive) |
| `-h` / `--help` | Show help |
| `-v` / `--version` | Show version |

## Examples

```bash
# Find who is using port 3000
PortKiller list | grep 3000

# Kill all Node processes
PortKiller kill -name node

# Kill the dev server without prompts
PortKiller kill -port 3000 -f
```
