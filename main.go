package main

import (
	"github.com/Awak3r/PortKiller/parser"
	"github.com/Awak3r/PortKiller/utils"
)

func main() {
	utils.EnsureRoot()
	p := utils.Collect()
	parser.ArgParse(p)
}
