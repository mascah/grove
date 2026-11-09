// Package cli parses commands and wires the other packages together. It is
// the only package that writes to the terminal.
package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

const usage = `usage: grove [command]

With no command, open the reader on the maintained specs and approved changes.

  init                       set up Grove in this repository
  instructions explore       print the instructions for exploring an idea
  instructions prototype "<question>"
                             print the instructions for answering a question with a
                             throwaway experiment
  instructions propose [--schema <name>] ["<description>"]
                             print the instructions for proposing a new change
  instructions propose <change-id>
                             print the instructions for revising a drafted change
  instructions plan ["<description>" | <change-id>]
                             print the instructions for planning a new cycle
  instructions plan <cycle-id>
                             print the instructions for revising a cycle
  instructions implement <change-id> | <cycle-id> [--unattended]
                             print the instructions for implementing an approved change,
                             or each change of a cycle, with the owner or asking no one
  instructions archive <change-id>
                             print the instructions for archiving a change
  claim [--retry] [<cycle-id>]
                             record an attempt at the cycle given, or else the oldest
                             marked unattended, on a new branch, and print its ID
  serve [--addr <host:port>] show the project's state on a local page, at
                             127.0.0.1:7070 unless another loopback address is given
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
	case len(args) == 3 && args[0] == "instructions" && args[1] == "prototype":
		err = instructionsPrototype(stdout, args[2])
	case len(args) == 2 && args[0] == "instructions" && args[1] == "propose":
		err = instructionsPropose(stdout, "", "")
	case len(args) == 3 && args[0] == "instructions" && args[1] == "propose":
		err = instructionsPropose(stdout, "", args[2])
	case len(args) >= 4 && len(args) <= 5 && args[0] == "instructions" && args[1] == "propose" && args[2] == "--schema":
		err = instructionsPropose(stdout, args[3], strings.Join(args[4:], ""))
	case len(args) == 2 && args[0] == "instructions" && args[1] == "plan":
		err = instructionsPlan(stdout, "")
	case len(args) == 3 && args[0] == "instructions" && args[1] == "plan":
		err = instructionsPlan(stdout, args[2])
	case len(args) == 3 && args[0] == "instructions" && args[1] == "implement":
		err = instructionsImplement(stdout, args[2], false)
	case len(args) == 4 && args[0] == "instructions" && args[1] == "implement" && args[3] == "--unattended":
		err = instructionsImplement(stdout, args[2], true)
	case len(args) == 3 && args[0] == "instructions" && args[1] == "archive":
		err = instructionsArchive(stdout, args[2])
	case len(args) > 0 && args[0] == "claim":
		f := flag.NewFlagSet("claim", flag.ContinueOnError)
		f.SetOutput(io.Discard)
		retry := f.Bool("retry", false, "")
		if f.Parse(args[1:]) != nil || f.NArg() > 1 {
			fmt.Fprint(stderr, usage)
			return 2
		}
		err = claim(stdout, f.Arg(0), *retry)
	case len(args) > 0 && args[0] == "serve":
		f := flag.NewFlagSet("serve", flag.ContinueOnError)
		f.SetOutput(io.Discard)
		addr := f.String("addr", "127.0.0.1:7070", "")
		if f.Parse(args[1:]) != nil || f.NArg() > 0 {
			fmt.Fprint(stderr, usage)
			return 2
		}
		err = serve(stdout, *addr)
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
