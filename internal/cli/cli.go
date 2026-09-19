// Package cli implements Grove's command interface.
package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/mascah/grove/internal/create"
	"github.com/mascah/grove/internal/project"
)

const usage = "Usage: grove [--project DIR] list | show ID [--json] | check | new TYPE TITLE [--slug SLUG]\n\n" +
	"  list       List records in the selected checkout\n" +
	"  show ID    Print the complete Markdown source for a record;\n" +
	"             --json prints {id, path, revision, source} instead\n" +
	"  check      Validate configuration, records, and relationships\n" +
	"  new        Create a work, question, or decision record with the next shared ID;\n" +
	"             put -- before a title that starts with a dash\n\n" +
	"--project DIR selects a directory containing grove.yaml.\n" +
	"Without it, search upward from the current directory, stopping at Git boundaries.\n" +
	"Project/file context is written to stderr; results are written to stdout.\n"

// Run returns 0 on success, 1 for inspection/output errors, and 2 for usage errors.
// cwd is explicit so callers and tests never need to change the process directory.
func Run(args []string, cwd string, out, errOut io.Writer) int {
	a, err := parseArgs(args)
	if err != nil {
		fmt.Fprintf(errOut, "grove: %s\n\n%s", err, usage)
		return 2
	}
	if a.help {
		return writeResult(out, errOut, []byte(usage))
	}
	p, ds := project.Load(cwd, a.project)
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
	switch a.command {
	case "new":
		path, err := create.New(p, a.kind, a.title, a.slug, time.Now(), errOut)
		if err != nil {
			for i, line := range strings.Split(err.Error(), "\n") {
				if i == 0 {
					line = "grove: " + line
				}
				fmt.Fprintln(errOut, visible(line))
			}
			return 1
		}
		return writeResult(out, errOut, []byte(path+"\n"))
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
			if r.ID == a.id {
				if _, err := fmt.Fprintf(errOut, "File: %s\n", visible(r.Path)); err != nil {
					return 1
				}
				if a.json {
					return writeResult(out, errOut, marshal(map[string]any{
						"id": r.ID, "path": r.Path, "revision": project.Revision(r.Source), "source": string(r.Source),
					}))
				}
				return writeResult(out, errOut, r.Source)
			}
		}
		fmt.Fprintf(errOut, "grove: record %s not found in this project\n", visible(a.id))
		return 1
	case "check":
		return writeResult(out, errOut, fmt.Appendf(nil, "OK: %d records\n", len(p.Records)))
	default:
		panic("validated command not handled")
	}
}

type invocation struct {
	project, command, id, kind, title, slug string
	help, json                              bool
}

// marshal encodes one flat object; the inputs are strings and booleans, which
// cannot fail to encode.
func marshal(object map[string]any) []byte {
	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	encoder.SetEscapeHTML(false)
	encoder.Encode(object)
	return buffer.Bytes()
}

func parseArgs(args []string) (a invocation, err error) {
	var positional []string
	literal := false
	// option consumes "--name VALUE" or "--name=VALUE" into *target, once.
	option := func(i *int, name, what string, target *string) (bool, error) {
		arg := args[*i]
		if arg != name && !strings.HasPrefix(arg, name+"=") {
			return false, nil
		}
		if *target != "" {
			return true, fmt.Errorf("%s may only be supplied once", name)
		}
		if arg == name {
			*i++
			if *i >= len(args) {
				return true, fmt.Errorf("%s requires a %s", name, what)
			}
			*target = args[*i]
		} else {
			*target = strings.TrimPrefix(arg, name+"=")
		}
		if strings.TrimSpace(*target) == "" {
			return true, fmt.Errorf("%s requires a nonempty %s", name, what)
		}
		return true, nil
	}
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if literal || !strings.HasPrefix(arg, "-") {
			positional = append(positional, arg)
			continue
		}
		if arg == "--" {
			literal = true
			continue
		}
		if arg == "--help" || arg == "-h" {
			a.help = true
			continue
		}
		if arg == "--json" {
			if a.json {
				return a, fmt.Errorf("--json may only be supplied once")
			}
			a.json = true
			continue
		}
		matched, err := option(&i, "--project", "directory", &a.project)
		if !matched && err == nil {
			matched, err = option(&i, "--slug", "slug", &a.slug)
		}
		if err != nil {
			return a, err
		}
		if !matched {
			return a, fmt.Errorf("unknown option %s", visible(arg))
		}
	}
	if a.help || (len(positional) == 1 && positional[0] == "help") {
		a.help = true
		return a, nil
	}
	if len(positional) == 0 {
		return a, fmt.Errorf("a command is required")
	}
	a.command = positional[0]
	if a.slug != "" && a.command != "new" {
		return a, fmt.Errorf("--slug applies only to new")
	}
	if a.json && a.command != "show" {
		return a, fmt.Errorf("--json applies only to show")
	}
	switch a.command {
	case "list", "check":
		if len(positional) != 1 {
			err = fmt.Errorf("%s takes no positional arguments", a.command)
		}
	case "show":
		if len(positional) != 2 {
			err = fmt.Errorf("show requires exactly one record ID")
		} else {
			a.id = positional[1]
		}
	case "new":
		if len(positional) != 3 {
			err = fmt.Errorf("new requires a record type and a title")
		} else {
			a.kind, a.title = positional[1], positional[2]
		}
	default:
		err = fmt.Errorf("unknown command %s", visible(a.command))
	}
	return a, err
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
