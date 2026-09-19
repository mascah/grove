// Package cli implements Grove's command interface.
package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/mascah/grove/internal/create"
	"github.com/mascah/grove/internal/project"
	"github.com/mascah/grove/internal/update"
	"github.com/mascah/grove/internal/versions"
)

const usage = "Usage: grove [--project DIR] list | show ID [--json] | check | new TYPE TITLE [--slug SLUG]\n" +
	"       grove [--project DIR] update ID --expect REVISION (--set FIELD=VALUE | --unset FIELD)...\n" +
	"       grove [--project DIR] versions [ID] [--json]\n" +
	"       grove [--project DIR] workspace --source SELECTOR [--json]\n\n" +
	"  list       List records in the selected checkout\n" +
	"  show ID    Print the complete Markdown source for a record;\n" +
	"             --json prints {id, path, revision, source} instead\n" +
	"  check      Validate configuration, records, and relationships\n" +
	"  new        Create a work, question, or decision record with the next shared ID;\n" +
	"             put -- before a title that starts with a dash\n" +
	"  update ID  Change frontmatter fields when the file still matches --expect\n" +
	"             (the revision from show --json); prints {id, path, revision, changed}.\n" +
	"             Lists are JSON arrays such as '[\"W-001\"]'; priority is 1-5.\n" +
	"  versions   Show each record's committed version on every local branch and live\n" +
	"             version in every worktree, with a selector per version; exit 1 if any\n" +
	"             source could not be inspected. Reads only; nothing is created.\n" +
	"  workspace  Print the project directory of the existing checkout holding the\n" +
	"             version selected by --source (a selector from versions), after\n" +
	"             checking it is still that version; --json adds checkout, record,\n" +
	"             branch, HEAD, and revision. Creates, switches, and edits nothing.\n\n" +
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
	// versions and workspace report this checkout as one live source among
	// others, so an invalid current project is attributed there rather than
	// ending the command.
	if len(ds) != 0 && (p == nil || (a.command != "versions" && a.command != "workspace")) {
		for _, d := range ds {
			fmt.Fprintln(errOut, visible(d.String()))
		}
		return 1
	}
	switch a.command {
	case "versions":
		return runVersions(p.Root, a, out, errOut)
	case "workspace":
		return runWorkspace(p.Root, a, out, errOut)
	case "update":
		res, err := update.Apply(p.Root, a.request, time.Now(), nil)
		if err != nil {
			report(errOut, err)
			return 1
		}
		result := marshal(map[string]any{"id": res.ID, "path": res.Path, "revision": res.Revision, "changed": res.Changed})
		if _, err := io.Copy(out, bytes.NewReader(result)); err != nil {
			state := "no change was needed for"
			if res.Changed {
				state = "the update was applied to"
			}
			fmt.Fprintf(errOut, "grove: write output: %s (%s %s; revision %s)\n", err, state, visible(res.Path), res.Revision)
			return 1
		}
		return 0
	case "new":
		path, err := create.New(p, a.kind, a.title, a.slug, time.Now(), errOut)
		if err != nil {
			report(errOut, err)
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
	project, command, id, kind, title, slug, source string
	help, json                                      bool
	request                                         update.Request
}

var revisionPattern = regexp.MustCompile(`^sha256:[0-9a-f]{64}$`)

// report writes a multi-line error with the grove: prefix on its first line.
func report(errOut io.Writer, err error) {
	for i, line := range strings.Split(err.Error(), "\n") {
		if i == 0 {
			line = "grove: " + line
		}
		fmt.Fprintln(errOut, visible(line))
	}
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
	fields := map[string]bool{}
	once := func(target *string) func(string) error {
		return func(value string) error {
			if *target != "" {
				return errors.New("may only be supplied once")
			}
			*target = value
			return nil
		}
	}
	field := func(name string) error {
		if fields[name] {
			return fmt.Errorf("mentions field %s more than once", visible(name))
		}
		fields[name] = true
		return nil
	}
	options := []struct {
		name, what string
		accept     func(string) error
	}{
		{"--project", "directory", once(&a.project)},
		{"--slug", "slug", once(&a.slug)},
		{"--source", "selector", once(&a.source)},
		{"--expect", "revision", once(&a.request.Expect)},
		{"--set", "FIELD=VALUE", func(value string) error {
			name, val, ok := strings.Cut(value, "=")
			if !ok || name == "" {
				return errors.New("requires FIELD=VALUE")
			}
			a.request.Set = append(a.request.Set, update.Field{Name: name, Value: val})
			return field(name)
		}},
		{"--unset", "field", func(value string) error {
			a.request.Unset = append(a.request.Unset, value)
			return field(value)
		}},
	}
	// option consumes "--name VALUE" or "--name=VALUE".
	option := func(i *int, name, what string, accept func(string) error) (bool, error) {
		arg := args[*i]
		if arg != name && !strings.HasPrefix(arg, name+"=") {
			return false, nil
		}
		var value string
		if arg == name {
			*i++
			if *i >= len(args) {
				return true, fmt.Errorf("%s requires a %s", name, what)
			}
			value = args[*i]
		} else {
			value = strings.TrimPrefix(arg, name+"=")
		}
		if strings.TrimSpace(value) == "" {
			return true, fmt.Errorf("%s requires a nonempty %s", name, what)
		}
		if err := accept(value); err != nil {
			return true, fmt.Errorf("%s %s", name, err)
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
		matched := false
		for _, o := range options {
			var err error
			if matched, err = option(&i, o.name, o.what, o.accept); err != nil {
				return a, err
			}
			if matched {
				break
			}
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
	if a.json && a.command != "show" && a.command != "versions" && a.command != "workspace" {
		return a, fmt.Errorf("--json applies only to show, versions, and workspace")
	}
	if a.source != "" && a.command != "workspace" {
		return a, fmt.Errorf("--source applies only to workspace")
	}
	if (a.request.Expect != "" || len(fields) != 0) && a.command != "update" {
		return a, fmt.Errorf("--expect, --set, and --unset apply only to update")
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
	case "versions":
		if len(positional) > 2 {
			err = fmt.Errorf("versions takes at most one record ID")
		} else if len(positional) == 2 {
			a.id = positional[1]
		}
	case "workspace":
		switch {
		case len(positional) != 1:
			err = fmt.Errorf("workspace takes no positional arguments")
		case a.source == "":
			err = fmt.Errorf("workspace requires --source SELECTOR from versions")
		default:
			_, err = versions.Parse(a.source)
		}
	case "new":
		if len(positional) != 3 {
			err = fmt.Errorf("new requires a record type and a title")
		} else {
			a.kind, a.title = positional[1], positional[2]
		}
	case "update":
		switch {
		case len(positional) != 2:
			err = fmt.Errorf("update requires exactly one record ID")
		case a.request.Expect == "":
			err = fmt.Errorf("update requires --expect REVISION from show --json")
		case !revisionPattern.MatchString(a.request.Expect):
			err = fmt.Errorf("--expect must be sha256: followed by 64 lowercase hexadecimal digits")
		case len(fields) == 0:
			err = fmt.Errorf("update requires at least one --set FIELD=VALUE or --unset FIELD")
		default:
			a.request.ID = positional[1]
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
