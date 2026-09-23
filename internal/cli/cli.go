// Package cli implements Grove's command interface.
package cli

import (
	"bytes"
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/mascah/grove"
	"github.com/mascah/grove/internal/attempt"
	"github.com/mascah/grove/internal/create"
	"github.com/mascah/grove/internal/handoff"
	"github.com/mascah/grove/internal/integrate"
	"github.com/mascah/grove/internal/project"
	"github.com/mascah/grove/internal/update"
	"github.com/mascah/grove/internal/versions"
)

const usage = "Usage: grove [--project DIR] [--json]\n" +
	"       grove [--project DIR] list [--status VALUE]... | show ID [--json] | brief [--json] | check\n" +
	"       grove [--project DIR] init\n" +
	"       grove guide work|shape | version\n" +
	"       grove [--project DIR] new TYPE TITLE [--slug SLUG]\n" +
	"       grove [--project DIR] update ID [--expect REVISION] (--set FIELD=VALUE | --unset FIELD)... [--commit]\n" +
	"       grove [--project DIR] approve ID VERDICT | feedback ID TEXT | integrate ID [--cleanup]\n" +
	"       grove [--project DIR] run ID --budget USD --permission-mode MODE [--model MODEL]\n" +
	"                                     [--branch NAME] [--worktree DIR]\n" +
	"       grove [--project DIR] attempts [ID] | attempt ATTEMPT [--json] | stop ATTEMPT\n" +
	"       grove [--project DIR] convert PATH --type TYPE --title TITLE [--slug SLUG]\n" +
	"       grove [--project DIR] versions [ID] [--json]\n" +
	"       grove [--project DIR] workspace --source SELECTOR [--json]\n" +
	"       grove [--project DIR] context WORK_ID... [--json] [--interaction interactive|headless]\n" +
	"                                     [--max-bytes N] [--include PATH]...\n\n" +
	"  (none)     Open the terminal board: work by its current state across every branch\n" +
	"             and checkout (or one checkout's own), each card's differing versions, and\n" +
	"             explicit selection of a version's existing workspace, printed like\n" +
	"             workspace (--json likewise).\n" +
	"             Needs a terminal on stdin and stderr; stdout may be redirected. Reads only.\n" +
	"  list       List records in the selected checkout; --status VALUE, repeatable, keeps\n" +
	"             only records in any given status (a value outside the vocabulary is refused)\n" +
	"  show ID    Print the complete Markdown source for a record;\n" +
	"             --json prints {id, path, revision, source} instead\n" +
	"  brief      Print the project brief that grove.yaml names with brief: PATH;\n" +
	"             --json prints {path, revision, source}. context never adds it by itself.\n" +
	"  check      Validate configuration, records, and relationships\n" +
	"  init       Set up the Git checkout at --project DIR (default: the current directory,\n" +
	"             which must be the checkout's top): grove.yaml, the record root, a\n" +
	"             placeholder brief, and the grove-work and grove-shape entrypoints for\n" +
	"             Claude Code and Codex. Existing files are kept; a file init wrote before\n" +
	"             (marked as managed) is updated when its template changed. Prints one line\n" +
	"             per path; on any conflict nothing is written and the reasons are printed.\n" +
	"  guide      Print the work or shaping guide this binary carries; the generated\n" +
	"             entrypoints read it from here, so the workflow version is the binary's.\n" +
	"  version    Print this binary's module version, its VCS revision when stamped, and\n" +
	"             a digest of the guides it carries.\n" +
	"  new        Create a work, question, decision, term, plan, review, or page record with\n" +
	"             the next shared ID, flat in the record root; a page is general knowledge\n" +
	"             with a title and no status.\n" +
	"             Put -- before a title that starts with a dash\n" +
	"  update ID  Change frontmatter fields; prints {id, path, revision, changed}. --expect\n" +
	"             REVISION (from show --json) refuses a file that no longer hashes to it, for a\n" +
	"             caller whose read may be old; omitted, the update applies to the file as it\n" +
	"             is. --commit then commits that one file with a generated message and adds\n" +
	"             commit to the result (null when nothing changed); other paths stay as they are.\n" +
	"             Lists are JSON arrays such as '[\"G-001\"]'; priority is 1-5. A plan or\n" +
	"             review names its work with work=[...]; a review's examined is a Git commit,\n" +
	"             as is work's candidate, required in review and, reachable from HEAD, for done.\n" +
	"             update accepts type=TYPE with whatever else the new type requires in the\n" +
	"             same update; the ID and path never change.\n" +
	"  approve    Record the owner's verdict on a work record in review, in a checkout of the\n" +
	"             branch that holds it: sets approved to the candidate, appends the verdict\n" +
	"             to the body, and commits that file alone. Refused where HEAD lacks the\n" +
	"             candidate, the record has uncommitted changes, or a commit after the\n" +
	"             candidate changed another file (that tip is a new candidate).\n" +
	"  feedback   Return a work record in review to active with the text appended to the body,\n" +
	"             committed alone in that same checkout; an approval is unset and the\n" +
	"             candidate kept, so earlier reviews still compare to it. Prints where to\n" +
	"             continue. Both print what update prints.\n" +
	"  integrate  Merge the one branch holding an approved candidate of ID into the target\n" +
	"             branch grove.yaml names, in that target's clean checkout, and mark ID done\n" +
	"             there, committed alone. Prints one line per fact as it holds: approval,\n" +
	"             merge (fast-forward or merge commit; a conflict is aborted and refused),\n" +
	"             done, and with --cleanup the worktree and branch removed, or kept with\n" +
	"             Git's reason. Every refusal comes before the merge; nothing undoes one.\n" +
	"  run        Start one bounded implementation attempt of a proposed or active work ID\n" +
	"             as a Grove-owned `claude -p \"/grove-work ID --interaction headless\"` process\n" +
	"             that outlives this terminal: in the branch's worktree (default worktree-ID\n" +
	"             under .claude/worktrees/, created from this checkout's HEAD or reused), with\n" +
	"             --max-budget-usd USD, --permission-mode MODE and --permission-prompts none, its\n" +
	"             raw output in files under the Git common directory. Refused while an attempt\n" +
	"             of ID runs or is orphaned, when ID is not proposed or active here or on its\n" +
	"             branch (a candidate in review awaits judgment), while an open question blocks\n" +
	"             ID in either place, when the record has uncommitted changes here, or when the\n" +
	"             worktree path is something else. Prints one line per\n" +
	"             fact and the attempt id. A result is facts, never acceptance: the record's own\n" +
	"             status on the branch is the handoff.\n" +
	"  attempts   List this repository's attempts, newest first, or those of one work ID:\n" +
	"             running (its owner holds the lock), finished (a result was written), orphaned\n" +
	"             (owner lost, provider alive) or interrupted (owner lost, nothing alive).\n" +
	"  attempt    Print one attempt's launch, status, event counts, result and file paths, and\n" +
	"             whether the record on the target changed since launch; --json prints the\n" +
	"             whole view. Reads files only: nothing is started or resumed.\n" +
	"  stop       End a running attempt through its owner (SIGINT ends the turn, SIGKILL after\n" +
	"             15 s) or an orphaned one directly; the result is written and the worktree and\n" +
	"             record are left as they are.\n" +
	"  convert    The one deliberate identity change. A Markdown document outside the record\n" +
	"             root becomes a new record with the document as its body and formerly: PATH;\n" +
	"             the original is left in place. Prints {from, from_path, id, path}. Bodies\n" +
	"             and Markdown links are never rewritten. A source already converted is\n" +
	"             refused, reserving nothing.\n" +
	"  versions   Show each record's committed version on every local branch and live\n" +
	"             version in every worktree, whether it is current or older and on the\n" +
	"             integration target, with a selector per version; exit 1 if any source\n" +
	"             could not be inspected. Reads only; nothing is created.\n" +
	"  workspace  Print the project directory of the existing checkout holding the\n" +
	"             version selected by --source (a selector from versions), after\n" +
	"             checking it is still that version; --json adds checkout, record,\n" +
	"             branch, HEAD, and revision. Creates, switches, and edits nothing.\n" +
	"  context    Assemble context for explicitly selected work in this checkout. Read in\n" +
	"             full, each with its exact revision: grove.yaml, the selected records, and\n" +
	"             every --include PATH (a required project-relative file). Listed, not read:\n" +
	"             prerequisites, blocking questions, plans and reviews whose work names a\n" +
	"             selected ID, related, member, and linked records with\n" +
	"             title, status, path, and revision, and the selected records' links with the\n" +
	"             path each resolves to (never opened, so not checked). Read a listed record\n" +
	"             with show ID; add a listed file, such as the current plan, with --include.\n" +
	"             The IDs are ordered prerequisites first, with Git identity. A missing,\n" +
	"             changed, or oversized source fails the command: nothing is truncated to\n" +
	"             fit --max-bytes (default 262144). --interaction records whether a person\n" +
	"             can answer (default interactive). Exit 0 means context was assembled, not\n" +
	"             that work is ready or authorized. Reads only.\n\n" +
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
	switch a.command {
	case "":
		return runBoard(a, cwd, out, errOut)
	case "version":
		return writeResult(out, errOut, []byte(versionLine()))
	case "guide":
		source, err := fs.ReadFile(grove.Guides, "docs/work-"+map[string]string{"work": "execution", "shape": "shaping"}[a.id]+".md")
		if err != nil {
			panic(err) // the two names were validated and both files are embedded
		}
		return writeResult(out, errOut, source)
	case "init":
		return runInit(cwd, a, out, errOut)
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
	case "context":
		return runContext(p.Root, a, out, errOut)
	case "update", "approve", "feedback":
		var res update.Result
		var err error
		switch a.command {
		case "update":
			res, err = update.Apply(p.Root, a.request, time.Now(), nil)
		case "approve":
			res, err = update.Approve(p.Root, a.id, a.title, time.Now())
		default:
			res, err = update.Feedback(p.Root, a.id, a.title, time.Now())
		}
		if err != nil {
			report(errOut, err)
			return 1
		}
		if a.command == "feedback" {
			// The actionable continuation: the work is active again here.
			branch, _ := update.Branch(p.Root)
			fmt.Fprintf(errOut, "Next: %s is active on branch %s in %s; continue there with /grove-work %s\n", res.ID, visible(cmp.Or(branch, "(detached HEAD)")), visible(p.Root), res.ID)
		}
		object := map[string]any{"id": res.ID, "path": res.Path, "revision": res.Revision, "changed": res.Changed}
		if a.request.Commit || a.command != "update" { // approve and feedback always commit
			object["commit"] = nil // a no-op commits nothing
			if res.Commit != "" {
				object["commit"] = res.Commit
			}
		}
		if _, err := io.Copy(out, bytes.NewReader(marshal(object))); err != nil {
			state := "no change was needed for"
			if res.Changed {
				state = "the update was applied to"
			}
			committed := ""
			if res.Commit != "" {
				committed = "; commit " + res.Commit
			}
			fmt.Fprintf(errOut, "grove: write output: %s (%s %s; revision %s%s)\n", err, state, visible(res.Path), res.Revision, committed)
			return 1
		}
		return 0
	case "integrate":
		err := integrate.Run(integrate.Request{Root: p.Root, ID: a.id, Cwd: cwd, Cleanup: a.cleanup}, time.Now(), func(fact string) {
			fmt.Fprintln(out, visible(fact))
		})
		if err != nil {
			report(errOut, err)
			return 1
		}
		return 0
	case "run":
		a.run.Root = p.Root
		l, err := attempt.Start(a.run, time.Now(), func(fact string) { fmt.Fprintln(out, visible(fact)) })
		if err != nil {
			report(errOut, err)
			return 1
		}
		fmt.Fprintf(out, "attempt: %s started; owner pid %d, session %s, budget %s USD, permission mode %s\n", l.Attempt, l.Owner, l.SessionID, l.BudgetUSD, l.PermissionMode)
		fmt.Fprintf(out, "inspect: grove attempt %s; stop: grove stop %s\n", l.Attempt, l.Attempt)
		return 0
	case "attempts":
		views, err := attempt.List(p.Root, a.id)
		if err != nil {
			report(errOut, err)
			return 1
		}
		return writeResult(out, errOut, attemptsTable(views))
	case "attempt":
		v, err := attempt.Show(p.Root, a.id)
		if err != nil {
			report(errOut, err)
			return 1
		}
		if a.json {
			data, _ := json.MarshalIndent(v, "", " ")
			return writeResult(out, errOut, append(data, '\n'))
		}
		return writeResult(out, errOut, []byte(attemptText(v)))
	case "stop":
		if err := attempt.Stop(p.Root, a.id, func(fact string) { fmt.Fprintln(out, visible(fact)) }); err != nil {
			report(errOut, err)
			return 1
		}
		return 0
	case "convert":
		c, err := update.Convert(p.Root, a.convert, errOut)
		code := 0
		if err != nil {
			report(errOut, err)
			if code = 1; c.ID == "" {
				return 1
			} // an incomplete conversion still prints the mapping it made
		}
		if _, err := io.Copy(out, bytes.NewReader(marshal(map[string]any{"from": c.From, "from_path": c.FromPath, "id": c.ID, "path": c.Path}))); err != nil {
			fmt.Fprintf(errOut, "grove: write output: %s (%s was converted to %s in %s)\n", err, visible(c.From), c.ID, visible(c.Path))
			return 1
		}
		return code
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
			if a.statuses != nil && !slices.Contains(a.statuses, r.Status) {
				continue
			}
			fmt.Fprintf(table, "%s\t%s\t%s\t%s\n", r.ID, r.Type, cmp.Or(r.Status, "-"), visible(r.Title)) // a page has no status
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
	case "brief":
		source, err := p.ReadBrief()
		if err != nil {
			report(errOut, err)
			return 1
		}
		if _, err := fmt.Fprintf(errOut, "File: %s\n", visible(p.Brief)); err != nil {
			return 1
		}
		if a.json {
			return writeResult(out, errOut, marshal(map[string]any{"path": p.Brief, "revision": project.Revision(source), "source": string(source)}))
		}
		return writeResult(out, errOut, source)
	case "check":
		return writeResult(out, errOut, fmt.Appendf(nil, "OK: %d records\n", len(p.Records)))
	default:
		panic("validated command not handled")
	}
}

type invocation struct {
	project, command, id, kind, title, slug, source string
	help, json, cleanup                             bool
	request                                         update.Request
	convert                                         update.ConvertRequest
	ids                                             []string // context
	statuses                                        []string // list
	options                                         handoff.Options
	run                                             attempt.Request
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
		{"--type", "record type", once(&a.convert.Type)},
		{"--title", "title", once(&a.convert.Title)},
		{"--source", "selector", once(&a.source)},
		{"--budget", "dollar amount", once(&a.run.BudgetUSD)},
		{"--permission-mode", "mode", once(&a.run.PermissionMode)},
		{"--model", "model", once(&a.run.Model)},
		{"--branch", "branch name", once(&a.run.Branch)},
		{"--worktree", "directory", once(&a.run.Worktree)},
		{"--expect", "revision", once(&a.request.Expect)},
		{"--interaction", "mode", func(value string) error {
			if value != "interactive" && value != "headless" {
				return errors.New("must be interactive or headless")
			}
			return once(&a.options.Interaction)(value)
		}},
		{"--max-bytes", "byte count", func(value string) error {
			if a.options.MaxBytes != 0 {
				return errors.New("may only be supplied once")
			}
			n, err := strconv.Atoi(value)
			if err != nil || n < 1 || n > handoff.LimitMaxBytes {
				return fmt.Errorf("must be an integer from 1 through %d", handoff.LimitMaxBytes)
			}
			a.options.MaxBytes = n
			return nil
		}},
		{"--status", "status", func(value string) error {
			if !slices.ContainsFunc(project.Types, func(t project.TypeInfo) bool { return slices.Contains(t.Statuses, value) }) {
				return fmt.Errorf("is not a status of any record type: %s", visible(value))
			}
			a.statuses = append(a.statuses, value)
			return nil
		}},
		{"--include", "project-relative path", func(value string) error {
			a.options.Include = append(a.options.Include, value)
			return nil
		}},
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
		if arg == "--commit" {
			if a.request.Commit {
				return a, fmt.Errorf("--commit may only be supplied once")
			}
			a.request.Commit = true
			continue
		}
		if arg == "--cleanup" {
			if a.cleanup {
				return a, fmt.Errorf("--cleanup may only be supplied once")
			}
			a.cleanup = true
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
	if len(positional) != 0 { // no command selects the board, under the same option rules
		a.command = positional[0]
	}
	if a.slug != "" && a.command != "new" && a.command != "convert" {
		return a, fmt.Errorf("--slug applies only to new and convert")
	}
	if (a.convert.Type != "" || a.convert.Title != "") && a.command != "convert" {
		return a, fmt.Errorf("--type and --title apply only to convert")
	}
	if a.json && a.command != "" && a.command != "show" && a.command != "brief" && a.command != "versions" && a.command != "workspace" && a.command != "context" && a.command != "attempt" {
		return a, fmt.Errorf("--json applies only to the board, show, brief, versions, workspace, context, and attempt")
	}
	if (a.run.BudgetUSD != "" || a.run.PermissionMode != "" || a.run.Model != "" || a.run.Branch != "" || a.run.Worktree != "") && a.command != "run" {
		return a, fmt.Errorf("--budget, --permission-mode, --model, --branch, and --worktree apply only to run")
	}
	if a.statuses != nil && a.command != "list" {
		return a, fmt.Errorf("--status applies only to list")
	}
	if a.source != "" && a.command != "workspace" {
		return a, fmt.Errorf("--source applies only to workspace")
	}
	if (a.options.Interaction != "" || a.options.MaxBytes != 0 || a.options.Include != nil) && a.command != "context" {
		return a, fmt.Errorf("--interaction, --max-bytes, and --include apply only to context")
	}
	if (a.request.Expect != "" || a.request.Commit || len(fields) != 0) && a.command != "update" {
		return a, fmt.Errorf("--expect, --set, --unset, and --commit apply only to update")
	}
	if a.cleanup && a.command != "integrate" {
		return a, fmt.Errorf("--cleanup applies only to integrate")
	}
	switch a.command {
	case "":
	case "list", "check", "brief", "init", "version":
		if len(positional) != 1 {
			err = fmt.Errorf("%s takes no positional arguments", a.command)
		}
	case "guide":
		if len(positional) != 2 || (positional[1] != "work" && positional[1] != "shape") {
			err = fmt.Errorf("guide requires one argument, work or shape")
		} else {
			a.id = positional[1]
		}
	case "show", "integrate":
		if len(positional) != 2 {
			err = fmt.Errorf("%s requires exactly one record ID", a.command)
		} else {
			a.id = positional[1]
		}
	case "run":
		switch {
		case len(positional) != 2:
			err = fmt.Errorf("run requires exactly one work ID")
		case a.run.BudgetUSD == "" || a.run.PermissionMode == "":
			err = fmt.Errorf("run requires --budget USD and --permission-mode MODE: Grove sets no default spend or permission profile")
		default:
			a.run.ID = positional[1]
		}
	case "attempts":
		if len(positional) > 2 {
			err = fmt.Errorf("attempts takes at most one work ID")
		} else if len(positional) == 2 {
			a.id = positional[1]
		}
	case "attempt", "stop":
		if len(positional) != 2 {
			err = fmt.Errorf("%s requires exactly one attempt id", a.command)
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
	case "context":
		if a.ids = positional[1:]; len(a.ids) == 0 {
			err = fmt.Errorf("context requires at least one work ID")
		}
	case "convert":
		if len(positional) != 2 {
			err = fmt.Errorf("convert requires exactly one document path")
		} else {
			a.convert.Source, a.convert.Slug = positional[1], a.slug
		}
	case "new":
		if len(positional) != 3 {
			err = fmt.Errorf("new requires a record type and a title")
		} else {
			a.kind, a.title = positional[1], positional[2]
		}
	case "approve", "feedback":
		what := map[string]string{"approve": "the verdict", "feedback": "the feedback text"}[a.command]
		if len(positional) != 3 {
			err = fmt.Errorf("%s requires a record ID and %s", a.command, what)
		} else {
			a.id, a.title = positional[1], positional[2]
		}
	case "update":
		switch {
		case len(positional) != 2:
			err = fmt.Errorf("update requires exactly one record ID")
		case a.request.Expect != "" && !revisionPattern.MatchString(a.request.Expect):
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
