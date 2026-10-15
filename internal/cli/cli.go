// Package cli parses commands and wires the other packages together. It is
// the only package that writes to the terminal.
package cli

import (
	"fmt"
	"os"
)

const usage = `usage: grove-next2 <command>

The commands are those of docs/walkthrough.md.
`

// Main runs one command and returns the exit status.
func Main(args []string) int {
	fmt.Fprint(os.Stderr, usage)
	return 2
}
