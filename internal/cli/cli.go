// Package cli parses commands and wires the other packages together. It is
// the only package that writes to the terminal.
package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/mascah/grove-next2/internal/doc"
)

const usage = `usage: grove-next2 <command>

  init                       set up Grove in this repository
  explore "<description>"    explore in a new workspace
  explore <id> [--resume]    a new conversation in that workspace, or the last one
  workspaces                 list the workspaces
  instructions propose <id>  print the propose instructions for that exploration
  check                      report problems in this working tree's Grove documents
`

// Main runs one command and returns the exit status.
func Main(args []string) int { return run(args, os.Stdout, os.Stderr) }

func run(args []string, stdout, stderr io.Writer) int {
	var err error
	switch {
	case len(args) == 1 && args[0] == "init":
		err = initProject(stdout)
	case len(args) == 2 && args[0] == "explore" && doc.ValidID(args[1]):
		err = reopen(stdout, args[1], false)
	case len(args) == 2 && args[0] == "explore":
		err = explore(stdout, args[1])
	case len(args) == 3 && args[0] == "explore" && args[2] == "--resume":
		err = reopen(stdout, args[1], true)
	case len(args) == 3 && args[0] == "instructions" && args[1] == "propose":
		err = instructionsPropose(stdout, args[2])
	case len(args) == 1 && args[0] == "check":
		err = check(stdout)
	case len(args) == 1 && args[0] == "workspaces":
		err = listWorkspaces(stdout)
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
