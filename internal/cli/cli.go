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
  instructions explore       print the instructions for exploring an idea
  instructions propose ["<description>" | <change-id>]
                             print the instructions for proposing a new change, or revising one
  instructions implement <change-id>
                             print the instructions for implementing an approved change
  review                     show what merging this branch into the target would approve
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
	case len(args) == 2 && args[0] == "instructions" && args[1] == "explore":
		err = instructionsExplore(stdout)
	case len(args) == 2 && args[0] == "instructions" && args[1] == "propose":
		err = instructionsPropose(stdout, "")
	case len(args) == 3 && args[0] == "instructions" && args[1] == "propose":
		err = instructionsPropose(stdout, args[2])
	case len(args) == 3 && args[0] == "instructions" && args[1] == "implement":
		err = instructionsImplement(stdout, args[2])
	case len(args) == 1 && args[0] == "review":
		err = review(stdout)
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
