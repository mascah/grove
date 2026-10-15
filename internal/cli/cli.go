// Package cli parses commands and wires the other packages together. It is
// the only package that writes to the terminal.
package cli

import (
	"fmt"
	"io"
	"os"
)

const usage = `usage: grove-next2 <command>

  init    set up Grove in this repository
`

// Main runs one command and returns the exit status.
func Main(args []string) int { return run(args, os.Stdout, os.Stderr) }

func run(args []string, stdout, stderr io.Writer) int {
	var err error
	switch {
	case len(args) == 1 && args[0] == "init":
		err = initProject(stdout)
	default:
		fmt.Fprint(stderr, usage)
		return 2
	}
	if err != nil {
		fmt.Fprintf(stderr, "grove-next2: %v\n", err)
		return 1
	}
	return 0
}
