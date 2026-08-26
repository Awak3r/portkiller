package commands

import "fmt"

func PrintUsage() {
  fmt.Println(`portkiller — kill processes by name or port

Usage:
  portkiller list [-name NAME | -port PORT]
  portkiller kill [-name NAME | -port PORT]
  portkiller -version
  portkiller -help`)
}