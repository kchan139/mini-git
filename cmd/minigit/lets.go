// Entry point and CLI interface

package main

import (
	"fmt"
	"os"

	"minigit/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
