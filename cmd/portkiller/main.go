package main

import (
	"log"

	"github.com/Awak3r/PortKiller/cli"
)

func main() {
	err := parser.ArgParse()
	if err != nil {
		log.Fatalf("%v", err)
	}
}
