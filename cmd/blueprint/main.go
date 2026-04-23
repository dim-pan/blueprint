package main

import (
	"os"

	"github.com/dim-pan/blueprint/src/cli"
)

func main() {
	os.Exit(cli.Run(os.Args[1:], os.Stdout, os.Stderr))
}
