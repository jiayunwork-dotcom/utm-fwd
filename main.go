package main

import (
	"os"

	"utm-fwd/internal/cli"
)

func main() {
	os.Exit(cli.Run(os.Args[1:]))
}
