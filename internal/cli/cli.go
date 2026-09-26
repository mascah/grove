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
	"github.com/mascah/grove/internal/sweep"
	"github.com/mascah/grove/internal/update"
	"github.com/mascah/grove/internal/versions"
)

const usage = "Usage: grove [--project DIR] [--json]\n" +
	"       grove [--project DIR] list [--status VALUE]... | show ID [--json] | brief [--json] | check\n" +
	"       grove [--project DIR] init [--check]\n" +
	"       grove guide work|shape|review|model [--entrypoint N] | version\n" +
	"       grove [--project DIR] new TYPE TITLE [--slug SLUG]\n" +
	"       grove [--project DIR] update ID [--expect REVISION] (--set FIELD=VALUE | --unset FIELD)... [--commit]\n" +
	"       grove [--project DIR] approve ID VERDICT | feedback ID TEXT | integrate ID [--cleanup]\n" +
	"       grove [--project DIR] resolve ID [--budget USD] [--permission-mode MODE] [--model MODEL]\n" +
	"                                     [--effort LEVEL]\n" +
	"       grove [--project DIR] sweep [--dry-run]\n" +
	"       grove [--project DIR] run WORK_ID... [--dry-run | --expect DIGEST] [--budget USD]\n" +
	"                                     [--permission-mode MODE] [--until plan] [--model MODEL]\n" +
	"                                     [--effort LEVEL] [--branch NAME] [--worktree DIR]\n" +
	"       grove [--project DIR] attempts [ID] | attempt ATTEMPT [--json] | stop ATTEMPT\n" +
	"       grove [--project DIR] convert PATH --type TYPE --title TITLE [--slug SLUG]\n" +
	"       grove [--project DIR] versions [ID] [--json]\n" +
	"       grove [--project DIR] workspace --source SELECTOR [--json]\n" +
	"       grove [--project DIR] context WORK_ID... [--json] [--interaction interactive|headless]\n" +
	"                                     [--max-bytes N] [--include PATH]...\n" +
	"       grove [--project DIR] deps [WORK_ID...] [--json]\n\n" +
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
	"             placeholder brief, the grove-work and grove-shape entrypoints for Claude\n" +
	"             Code and Codex, and Claude Code's grove-reviewer agent definition, which the\n" +
	"             work guide reviews through. Existing files are kept; a file init wrote before\n" +
	"             (marked as managed) is updated when its template changed. Prints one line\n" +
	"             per path; on any conflict nothing is written and the reasons are printed.\n" +
	"             --check writes nothing and prints each entrypoint as current, compatible,\n" +
	"             legacy (no entrypoint revision), incompatible (a revision this binary does\n" +
	"             not serve), missing, custom (unmarked, not judged) or conflict; exit 1 if\n" +
	"             any is legacy, incompatible, missing or a conflict.\n" +
	"  guide      Print the work, shaping or review guide, or the record model they cite,\n" +
	"             that this binary carries; the generated entrypoints read the guides from\n" +
	"             here, so the workflow version is the binary's. --entrypoint N is how an\n" +
	"             entrypoint of revision N asks; a revision this binary does not serve is\n" +
	"             refused.\n" +
	"  version    Print this binary's version and commit, and digests of the guides and\n" +
	"             record model and of all the content it ships.\n" +
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
	"             candidate changed another file (that tip is a new candidate). Work records\n" +
	"             whose candidate is the same commit are one group, handed off together from\n" +
	"             one selection: each is approved on its own, and their record files do not\n" +
	"             count as later changes.\n" +
	"  feedback   Return a work record in review to active with the text appended to the body,\n" +
	"             committed alone in that same checkout; an approval is unset and the\n" +
	"             candidate kept, so earlier reviews still compare to it. Every other member\n" +
	"             of its group in review is reopened the same way, each committed alone, with\n" +
	"             a line naming this feedback. Prints where to continue. Both print what\n" +
	"             update prints.\n" +
	"  integrate  Merge the one branch holding an approved candidate of ID into the target\n" +
	"             branch grove.yaml names, in that target's clean checkout, and mark ID done\n" +
	"             there, committed alone. A candidate shared by a group integrates the group:\n" +
	"             refused unless every member is approved in review, then one merge and done\n" +
	"             for each member, committed alone. Prints one line per fact as it holds: approval,\n" +
	"             merge (fast-forward or merge commit; a conflict is aborted and refused),\n" +
	"             done, and with --cleanup the worktree and branch removed, or kept with\n" +
	"             Git's reason. Every refusal comes before the merge; nothing undoes one.\n" +
	"  resolve    For a work record in review whose candidate conflicts with the target, from\n" +
	"             any checkout: record feedback naming the target commit and the conflicting\n" +
	"             files, committed in the branch's checkout, and start one attempt there, as run\n" +
	"             would, to merge that commit, resolve those files, verify, and hand off the\n" +
	"             merge as a new candidate. The records sharing the candidate go with it. Refused,\n" +
	"             with nothing written, without a target or a conflict, or while an attempt of\n" +
	"             the work runs; the launch flags and their run: defaults are run's.\n" +
	"  sweep      Act on every candidate in review under the standing policy: grove.yaml's\n" +
	"             policy:, committed, in the target's checkout. A conflict gets one resolution\n" +
	"             attempt, as resolve starts, once per target commit and within the policy's\n" +
	"             budget; a clean candidate whose current review ends \"Open findings: none\",\n" +
	"             with no never path and within max_lines, is merged with the target in a\n" +
	"             temporary worktree, verified there with the policy's commands, then approved\n" +
	"             and, with integrate: true, integrated, each attributed to the policy's\n" +
	"             revision. Anything else waits for the owner, with the reason. Prints one line\n" +
	"             per candidate and per act; --dry-run prints what would happen to each and\n" +
	"             why, and writes nothing. Refused without a policy: nothing is automatic.\n" +
	"  run        Start one bounded implementation attempt of an explicit selection of\n" +
	"             proposed or active work as one Grove-owned\n" +
	"             `claude -p \"/grove-work ID... --interaction headless\"` process that outlives\n" +
	"             this terminal, the IDs passed as given: in the branch's worktree (default\n" +
	"             worktree- plus the IDs joined by - under .claude/worktrees/, created from this\n" +
	"             checkout's HEAD or reused), with --max-budget-usd USD over the whole selection,\n" +
	"             --permission-mode MODE and --permission-prompts none, its raw output in files\n" +
	"             under the Git common directory. Members are implemented one at a time in\n" +
	"             dependency order; nothing outside the selection is added, and the complete\n" +
	"             members are handed off together on one shared candidate. --until plan ends the\n" +
	"             attempt at committed plans, leaving statuses as found; --model and --effort are\n" +
	"             passed to the provider. --budget and --permission-mode are required unless\n" +
	"             grove.yaml's run: sets them; it may set --model and --effort too, and a flag\n" +
	"             overrides it. --dry-run checks and prints the assignment without writing or\n" +
	"             starting anything: order, each member's revision and whether it can start or\n" +
	"             waits (an open question, a prerequisite outside the selection the base does not\n" +
	"             hold, or a selected one that waits), the outside prerequisites, the base,\n" +
	"             bounds, review boundary, continuation policy, and a digest; --expect DIGEST\n" +
	"             refuses a launch whose assignment no longer has that digest. What ran is\n" +
	"             recorded with the digest of the worktree's grove-reviewer definition, or its\n" +
	"             absence, which is warned of. Refused when the worktree would not hold the\n" +
	"             committed grove-work skill (.claude/skills/grove-work/SKILL.md, which init\n" +
	"             writes), when that skill or the reviewer is marked with an entrypoint revision\n" +
	"             this binary does not serve (or none), while an attempt whose selection shares a\n" +
	"             member runs or is orphaned, when a member is not proposed or active here or\n" +
	"             on its branch (a candidate in review awaits judgment), when no member can\n" +
	"             start (with --until plan only an open question stops one), when a member\n" +
	"             record has uncommitted changes here, or when the worktree path is something\n" +
	"             else. Prints one line per fact and the attempt id. A result is facts, never\n" +
	"             acceptance: the records' own statuses on the branch are the handoff.\n" +
	"  attempts   List this repository's attempts, newest first, or those whose selection\n" +
	"             includes one work ID:\n" +
	"             running (its owner holds the lock), finished (a result was written), orphaned\n" +
	"             (owner lost, provider alive) or interrupted (owner lost, nothing alive).\n" +
	"  attempt    Print one attempt's launch, status, event counts, result and file paths,\n" +
	"             each member's state on the branch (awaiting judgment, active, not started, or\n" +
	"             waiting and on what), and whether a member record on the target changed since\n" +
	"             launch; --json prints the\n" +
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
	"             that work is ready or authorized. Reads only.\n" +
	"  deps       Show how unfinished work in this checkout depends on other work: one row\n" +
	"             each with its group, layer, prerequisites and the work it unlocks, then\n" +
	"             the prerequisites that are not unfinished. With work IDs, preview that\n" +
	"             selection: its order (as context gives it), every prerequisite outside it,\n" +
	"             listed and never added, and the open questions blocking any of them. Each\n" +
	"             candidate's delivery is Git ancestry into HEAD and the target; notes say\n" +
	"             where other branches and checkouts hold other versions. Reads only; exit 1\n" +
	"             if any source could not be inspected.\n\n" +
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
		return writeResult(out, errOut, []byte(grove.Identity().String()+"\n"))
	case "guide":
		if a.entrypoint != "" && !grove.SupportsEntrypoint(a.entrypoint) {
			report(errOut, fmt.Errorf("the entrypoint that asked for this guide is revision %s, and this grove serves %s; nothing was printed.\n"+
				"Rerun `grove init` with this grove to rewrite the entrypoints it manages (`grove init --check` lists them),\n"+
				"commit them, and start a new session; or run the grove that wrote them.", visible(a.entrypoint), grove.ServedEntrypoints()))
			return 1
		}
		source, err := fs.ReadFile(grove.Guides, grove.GuideFiles[a.id])
		if err != nil {
			panic(err) // the name was validated and every file is embedded
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
	case "deps":
		return runDeps(p, a, out, errOut)
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
		object := map[string]any{"id": res.ID, "path": res.Path, "revision": res.Revision, "changed": res.Changed}
		if a.command == "feedback" {
			// The actionable continuation: the work, and any group it shared a
			// candidate with, is active again here.
			ids := []string{res.ID}
			var reopened []map[string]any
			for _, o := range res.Reopened {
				fmt.Fprintf(errOut, "Reopened: %s shared the candidate and is active again, commit %s\n", o.ID, o.Commit)
				ids = append(ids, o.ID)
				reopened = append(reopened, map[string]any{"id": o.ID, "path": o.Path, "revision": o.Revision, "commit": o.Commit})
			}
			if reopened != nil {
				object["reopened"] = reopened
			}
			branch, _ := update.Branch(p.Root)
			fmt.Fprintf(errOut, "Next: %s active on branch %s in %s; continue there with /grove-work %s\n", strings.Join(ids, ", ")+map[bool]string{true: " is", false: " are"}[len(ids) == 1], visible(cmp.Or(branch, "(detached HEAD)")), visible(p.Root), strings.Join(ids, " "))
			if len(ids) > 1 && branch != "" {
				fmt.Fprintf(errOut, "or launch the group again: grove run %s --branch %s\n", strings.Join(ids, " "), visible(branch))
			}
		}
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
		// Without grove.yaml's run: defaults the flags are required, which
		// only the loaded project can tell, so it is still a usage error.
		if d := attempt.Defaulted(a.run, p.Run); d.BudgetUSD == "" || d.PermissionMode == "" {
			fmt.Fprintf(errOut, "grove: %s\n\n%s", attempt.ErrUnsupplied, usage)
			return 2
		}
		a.run.Root = p.Root
		if a.dryRun {
			l, err := attempt.Preview(a.run)
			if err != nil {
				report(errOut, err)
				return 1
			}
			return writeResult(out, errOut, []byte(previewText(l)))
		}
		l, err := attempt.Start(a.run, time.Now(), func(fact string) { fmt.Fprintln(out, visible(fact)) })
		if err != nil {
			report(errOut, err)
			return 1
		}
		fmt.Fprintf(out, "attempt: %s started; owner pid %d, session %s, budget %s USD, permission mode %s\n", l.Attempt, l.Owner, l.SessionID, l.BudgetUSD, l.PermissionMode)
		fmt.Fprintf(out, "requested: %s\n", visible(attempt.Requested(l)))
		if len(l.Members()) > 1 { // what the preview showed, as launched
			for _, fact := range attempt.Explain(l, visible) {
				fmt.Fprintln(out, fact)
			}
		}
		fmt.Fprintf(out, "inspect: grove attempt %s; stop: grove stop %s\n", l.Attempt, l.Attempt)
		return 0
	case "resolve":
		a.run.Root, a.run.IDs = p.Root, []string{a.id}
		l, err := attempt.Resolve(a.run, nil, time.Now(), func(fact string) { fmt.Fprintln(out, visible(fact)) })
		if err != nil {
			report(errOut, err)
			return 1
		}
		fmt.Fprintf(out, "attempt: %s started on %s; owner pid %d, budget %s USD, permission mode %s\n", l.Attempt, visible(l.Branch), l.Owner, l.BudgetUSD, visible(l.PermissionMode))
		fmt.Fprintf(out, "inspect: grove attempt %s; stop: grove stop %s\n", l.Attempt, l.Attempt)
		return 0
	case "sweep":
		s, err := sweep.Plan(p.Root)
		if err != nil {
			report(errOut, err)
			return 1
		}
		if len(s.Items) == 0 {
			fmt.Fprintln(out, "no candidate is in review")
		}
		if a.dryRun {
			for _, it := range s.Items {
				fmt.Fprintf(out, "%s on %s: %s: %s\n", it.ID, visible(it.Branch), it.Act, visible(it.Why))
			}
			return 0
		}
		s.Run(time.Now(), func(fact string) { fmt.Fprintln(out, visible(fact)) })
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
	entrypoint                                      string // guide
	help, json, cleanup, check, dryRun              bool
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
		{"--expect", "revision", once(&a.request.Expect)},
		{"--entrypoint", "revision", once(&a.entrypoint)},
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
		if arg == "--dry-run" {
			if a.dryRun {
				return a, fmt.Errorf("--dry-run may only be supplied once")
			}
			a.dryRun = true
			continue
		}
		if arg == "--check" {
			if a.check {
				return a, fmt.Errorf("--check may only be supplied once")
			}
			a.check = true
			continue
		}
		if arg == "--cleanup" {
			if a.cleanup {
				return a, fmt.Errorf("--cleanup may only be supplied once")
			}
			a.cleanup = true
			continue
		}
		matched, err := attempt.Flag(args, &i, &a.run)
		if err != nil {
			return a, err
		}
		for _, o := range options {
			if matched {
				break
			}
			if matched, err = option(&i, o.name, o.what, o.accept); err != nil {
				return a, err
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
	if a.json && a.command != "" && a.command != "show" && a.command != "brief" && a.command != "versions" && a.command != "workspace" && a.command != "context" && a.command != "deps" && a.command != "attempt" {
		return a, fmt.Errorf("--json applies only to the board, show, brief, versions, workspace, context, deps, and attempt")
	}
	if (a.run.BudgetUSD != "" || a.run.PermissionMode != "" || a.run.Model != "" || a.run.Effort != "") && a.command != "run" && a.command != "resolve" {
		return a, fmt.Errorf("--budget, --permission-mode, --model, and --effort apply only to run and resolve")
	}
	if (a.run.Until != "" || a.run.Branch != "" || a.run.Worktree != "") && a.command != "run" {
		return a, fmt.Errorf("--until, --branch, and --worktree apply only to run")
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
	if a.request.Expect != "" && a.command != "update" && a.command != "run" {
		return a, fmt.Errorf("--expect applies only to update and run")
	}
	if (a.request.Commit || len(fields) != 0) && a.command != "update" {
		return a, fmt.Errorf("--set, --unset, and --commit apply only to update")
	}
	if a.dryRun && a.command != "run" && a.command != "sweep" {
		return a, fmt.Errorf("--dry-run applies only to run and sweep")
	}
	if a.cleanup && a.command != "integrate" {
		return a, fmt.Errorf("--cleanup applies only to integrate")
	}
	if a.check && a.command != "init" {
		return a, fmt.Errorf("--check applies only to init")
	}
	if a.entrypoint != "" && a.command != "guide" {
		return a, fmt.Errorf("--entrypoint applies only to guide")
	}
	switch a.command {
	case "":
	case "list", "check", "brief", "init", "version", "sweep":
		if len(positional) != 1 {
			err = fmt.Errorf("%s takes no positional arguments", a.command)
		}
	case "guide":
		if len(positional) != 2 || grove.GuideFiles[positional[1]] == "" {
			err = fmt.Errorf("guide requires one argument, work, shape, review or model")
		} else {
			a.id = positional[1]
		}
	case "show", "integrate", "resolve":
		if len(positional) != 2 {
			err = fmt.Errorf("%s requires exactly one record ID", a.command)
		} else {
			a.id = positional[1]
		}
	case "run":
		switch {
		case len(positional) < 2:
			err = fmt.Errorf("run requires at least one work ID")
		case a.request.Expect != "" && a.dryRun:
			err = fmt.Errorf("--expect checks a launch against a --dry-run's digest; give one or the other")
		case a.request.Expect != "" && !revisionPattern.MatchString(a.request.Expect):
			err = fmt.Errorf("--expect must be the digest --dry-run printed: sha256: followed by 64 lowercase hexadecimal digits")
		default:
			a.run.IDs, a.run.Digest, a.request.Expect = positional[1:], a.request.Expect, ""
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
	case "deps":
		a.ids = positional[1:]
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
