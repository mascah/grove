// Package cli parses commands and wires the other packages together. It is
// the only package that writes to the terminal.
package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mascah/grove/internal/instruct"
)

const usage = `usage: grove [command]

With no command, open the reader on the maintained specs and approved changes.

  init                       set up Grove in this repository
  instructions explore       print the instructions for exploring an idea
  instructions propose [--schema <name>] ["<description>"]
                             print the instructions for proposing a new change
  instructions propose <change-id>
                             print the instructions for revising a drafted change
  instructions implement <change-id>
                             print the instructions for implementing an approved change
                             with the owner
  instructions apply [<change-id>]
                             print the instructions for applying an approved change,
                             or every approved change with none
  instructions archive <change-id>
                             print the instructions for archiving a change
  run [--retry] [--timeout <duration>]
                             apply every approved change in one Claude Code session
                             on a new branch, then pause until they are archived
  archive <change-id>        move a change to changes/archive/
  review                     show what merging this branch into the target would approve
  specs                      list the maintained specs
  changes [--json]           list the approved changes, or print them as JSON
  check                      report problems in this working tree's Grove documents
  version                    print Grove's version
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
		err = instructionsPropose(stdout, "", "")
	case len(args) == 3 && args[0] == "instructions" && args[1] == "propose":
		err = instructionsPropose(stdout, "", args[2])
	case len(args) >= 4 && len(args) <= 5 && args[0] == "instructions" && args[1] == "propose" && args[2] == "--schema":
		err = instructionsPropose(stdout, args[3], strings.Join(args[4:], ""))
	case len(args) == 2 && args[0] == "instructions" && args[1] == "apply":
		err = instructionsApplyAll(stdout)
	case len(args) == 3 && args[0] == "instructions" && args[1] == "apply":
		err = instructionsApply(stdout, args[2], instruct.Apply)
	case len(args) == 3 && args[0] == "instructions" && args[1] == "implement":
		err = instructionsApply(stdout, args[2], instruct.Implement)
	case len(args) == 3 && args[0] == "instructions" && args[1] == "archive":
		err = instructionsArchive(stdout, args[2])
	case len(args) > 0 && args[0] == "run":
		f := flag.NewFlagSet("run", flag.ContinueOnError)
		f.SetOutput(io.Discard)
		retry, timeout := f.Bool("retry", false, ""), f.Duration("timeout", defaultTimeout, "")
		if f.Parse(args[1:]) != nil || f.NArg() > 0 {
			fmt.Fprint(stderr, usage)
			return 2
		}
		err = runAttempt(stdout, *retry, *timeout)
	case len(args) == 2 && args[0] == "archive":
		err = archive(stdout, args[1])
	case len(args) == 1 && args[0] == "review":
		err = review(stdout)
	case len(args) == 1 && args[0] == "changes":
		err = listChanges(stdout, false)
	case len(args) == 2 && args[0] == "changes" && args[1] == "--json":
		err = listChanges(stdout, true)
	case len(args) == 1 && args[0] == "check":
		err = check(stdout)
	case len(args) == 1 && args[0] == "version":
		printVersion(stdout)
	default:
		fmt.Fprint(stderr, usage)
		return 2
	}
	if err != nil {
		fmt.Fprintf(stderr, "grove: %v\n", err)
		return 1
	}
	return 0
}
