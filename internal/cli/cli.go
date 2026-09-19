// Package cli implements Grove's read-only command interface.
package cli

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/mascah/grove/internal/project"
)

const usage = "Usage: grove [--project DIR] list | show ID | check\n\n" +
	"  list       List records in the selected checkout\n" +
	"  show ID    Print the complete Markdown source for a record\n" +
	"  check      Validate configuration, records, and relationships\n\n" +
	"--project DIR selects a directory containing grove.yaml.\n" +
	"Without it, search upward from the current directory, stopping at Git boundaries.\n" +
	"Project/file context is written to stderr; results are written to stdout.\n"

// Run returns 0 on success, 1 for inspection/output errors, and 2 for usage errors.
// cwd is explicit so callers and tests never need to change the process directory.
func Run(args []string, cwd string, out, errOut io.Writer) int {
	selected, command, id, help, err := parseArgs(args)
	if err != nil {
		fmt.Fprintf(errOut, "grove: %s\n\n%s", err, usage)
		return 2
	}
	if help {
		return writeResult(out, errOut, []byte(usage))
	}
	p, ds := project.Load(cwd, selected)
	if p != nil {
		if _, err := fmt.Fprintf(errOut, "Project: %s\n", visible(p.Root)); err != nil {
			return 1
		}
	}
	if len(ds) != 0 {
		for _, d := range ds {
			fmt.Fprintln(errOut, visible(d.String()))
		}
		return 1
	}
	switch command {
	case "list":
		var buffer bytes.Buffer
		table := tabwriter.NewWriter(&buffer, 0, 4, 2, ' ', 0)
		fmt.Fprintln(table, "ID\tTYPE\tSTATUS\tTITLE")
		for _, r := range p.Records {
			fmt.Fprintf(table, "%s\t%s\t%s\t%s\n", r.ID, r.Type, r.Status, visible(r.Title))
		}
		table.Flush() // The destination is a bytes.Buffer, whose writes cannot fail.
		return writeResult(out, errOut, buffer.Bytes())
	case "show":
		for _, r := range p.Records {
			if r.ID == id {
				if _, err := fmt.Fprintf(errOut, "File: %s\n", visible(r.Path)); err != nil {
					return 1
				}
				return writeResult(out, errOut, r.Source)
			}
		}
		fmt.Fprintf(errOut, "grove: record %s not found in this project\n", visible(id))
		return 1
	case "check":
		return writeResult(out, errOut, fmt.Appendf(nil, "OK: %d records\n", len(p.Records)))
	default:
		panic("validated command not handled")
	}
}

func parseArgs(args []string) (selected, command, id string, help bool, err error) {
	var positional []string
	projectSet, literal := false, false
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case !literal && arg == "--":
			literal = true
		case !literal && (arg == "--help" || arg == "-h"):
			help = true
		case !literal && (arg == "--project" || strings.HasPrefix(arg, "--project=")):
			if projectSet {
				err = fmt.Errorf("--project may only be supplied once")
				return
			}
			projectSet = true
			if arg == "--project" {
				i++
				if i >= len(args) {
					err = fmt.Errorf("--project requires a directory")
					return
				}
				selected = args[i]
			} else {
				selected = strings.TrimPrefix(arg, "--project=")
			}
			if strings.TrimSpace(selected) == "" {
				err = fmt.Errorf("--project requires a nonempty directory")
				return
			}
		case !literal && strings.HasPrefix(arg, "-"):
			err = fmt.Errorf("unknown option %s", visible(arg))
			return
		default:
			positional = append(positional, arg)
		}
	}
	if help || (len(positional) == 1 && positional[0] == "help") {
		help = true
		return
	}
	if len(positional) == 0 {
		err = fmt.Errorf("a command is required")
		return
	}
	command = positional[0]
	switch command {
	case "list", "check":
		if len(positional) != 1 {
			err = fmt.Errorf("%s takes no positional arguments", command)
		}
	case "show":
		if len(positional) != 2 {
			err = fmt.Errorf("show requires exactly one record ID")
		} else {
			id = positional[1]
		}
	default:
		err = fmt.Errorf("unknown command %s", visible(command))
	}
	return
}

func writeResult(out, errOut io.Writer, content []byte) int {
	if _, err := io.Copy(out, bytes.NewReader(content)); err != nil {
		fmt.Fprintf(errOut, "grove: write output: %s\n", err)
		return 1
	}
	return 0
}

// Escape control characters in one-line output while retaining readable Unicode.
// show intentionally emits the original source instead.
func visible(value string) string {
	quoted := strconv.Quote(value)
	return quoted[1 : len(quoted)-1]
}
