// Package cli parses commands and wires the other packages together. It is
// the only package that writes to the terminal.
package cli

import (
	"fmt"
	"io"
	"os"
)

const usage = `usage: grove-next2 [command]

With no command, open the reader on the maintained specs and approved changes.

  init                       set up Grove in this repository
  instructions implement <change-id>
                             print the implement instructions for an approved change
  specs                      list the maintained specs
  changes                    list the approved changes
  check                      report problems in this working tree's Grove documents
`

// Main runs one command and returns the exit status.
func Main(args []string) int { return run(args, os.Stdout, os.Stderr) }

func run(args []string, stdout, stderr io.Writer) int {
	var err error
	switch {
	case len(args) == 0:
		err = reader()
	case len(args) == 1 && args[0] == "specs":
		err = listSpecs(stdout)
	case len(args) == 1 && args[0] == "init":
		err = initProject(stdout)
	case len(args) == 3 && args[0] == "instructions" && args[1] == "implement":
		err = instructionsImplement(stdout, args[2])
	case len(args) == 1 && args[0] == "changes":
		err = listChanges(stdout)
	case len(args) == 1 && args[0] == "check":
		err = check(stdout)
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
