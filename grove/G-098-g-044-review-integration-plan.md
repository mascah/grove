---
id: "G-098"
type: plan
title: "G-044 review actions and local integration plan"
status: current
created: "2026-09-23T00:17:01Z"
updated: "2026-09-23T00:17:10Z"
work: ["G-044"]
---

## Inputs

Prepared on 2026-09-22 for [G-044](G-044-review-integration.md) at revision
`493bf1c6`, on branch `worktree-G-044` from `main` `c4b3aea`, which holds
both prerequisites merged: G-038 (candidate `fae1e4c`) and G-043 (candidate
`178b305`). Read in full: G-044, G-038 with its plan
[G-073](G-073-review-lifecycle-plan.md), G-043 with its plan
[G-096](G-096-g-043-board-and-detail-design-vi.md) (whose "review view"
section this plan builds), G-035, G-046, G-064, the terms G-057 to G-060, the
brief, the record model, the work guide's "Judging and integrating" section,
the G-039 trial evidence [G-078](G-078-g-039-trial-evidence-for-the-in.md),
and the code: `internal/update`, `internal/cli`, `internal/repo`,
`internal/tui` (`model.go`, `detail.go`, `view.go`, `search.go`, `run.go`),
`internal/versions` (`versions.go`, `current.go`), the justfile's
`clean-merged` recipe, and the merges on `main`.

What the loop lacks today, from G-078 and the guide: approval "has no CLI
form: it is either told to an agent or a merge, a hand edit and an `update`
run by hand"; the verdict was not quoted when G-076 was closed; `update`
cannot tell the target from the work branch; the board shows a review's
`examined` against the candidate but offers no action; and nothing in Grove
merges, cleans up, or shows a candidate's changed files.

## Design

The four choices G-044 asks to resolve before building, each with the
alternative not taken. They are proposed here and put to the owner (see
Open choices); everything after them is a routine technical choice inside
G-044's outcome.

### 1. Approval representation: an `approved` field bound to the candidate

Work gains one optional field, `approved`: a quoted Git commit in the same
form as `candidate`. The record is valid only when `approved` names the
`candidate` exactly and the status is `review` or `done`. So:

- Approval is a fact software acts on (G-064), not prose. `grove approve`
  writes it; `integrate` requires it; `check` verifies the binding.
- A changed candidate cannot inherit approval: `update --set candidate=NEW`
  on an approved record is refused unless the same call unsets `approved`,
  and a commit after the candidate on the branch is caught by `integrate`
  (below), since the tip then differs from the candidate in more than the
  record file. Feedback that returns work to `active` must unset
  `approved` in the same update, so nothing survives silently.
- The verdict stays where G-059 puts it, quoted in the record: `approve`
  appends one paragraph to the end of the body, `Verdict on candidate
  178b305, 2026-09-22: <the owner's words>`, and commits the record alone.
  Feedback appends `Feedback on candidate 178b305, 2026-09-22: <text>` the
  same way. A body append is the one body edit this adds: it changes no
  existing byte, needs no heading parsing, and lands in or after the Next
  section, where the guide already puts feedback.
- Attribution is the commit: `approve`, `feedback` and `integrate` always
  commit the record's file alone with the checkout's Git identity, as
  `update --commit` does. There is no author field.

Not taken: verdict only at integration (one command that merges and closes);
it leaves nothing recorded when the merge is refused, and cannot show a
stale approval. Also not taken: an approval record type; a review record is
evidence and G-038 chose no disposition field.

### 2. Merge strategy: a plain merge into the target's checkout

`integrate` runs `git merge --no-edit BRANCH` in a checkout whose HEAD is
the configured `target`: a fast-forward when `main` has not moved, a merge
commit otherwise, which is how every integration on this `main` was done
(`Merge branch 'worktree-G-040'`). The candidate keeps its identity through
either, so the merged record's `candidate` still names the commit that was
reviewed and approved, and `done`'s ancestry check holds. A conflict is
aborted (`git merge --abort`) and reported with Git's words; the target is
left as it was and nothing is written to any record. No squash and no
rebase: both make a different commit from the candidate, which approval was
of.

Before merging, `integrate` checks, and refuses with a named reason:
`grove.yaml` names no `target`; this checkout is not on the target branch;
the target checkout has staged or unstaged changes to tracked files
(untracked files are left alone: Git refuses the merge itself if one is in
the way); no branch holds the record in `review` with `approved` equal to
its candidate (an unapproved one is named with the `approve` command to
run); more than one branch does (both named); or the branch tip differs
from the candidate in any file but the record's (a commit after the
handoff is a new candidate: the tip is named and approval must be given
again). Every refusal happens before the merge, so a refused integration
changes nothing.

### 3. Post-merge status: `integrate` writes done on the target

After the merge succeeds, `integrate` runs the existing update path in the
target's checkout: `status=done`, the candidate unchanged, committed alone
(`docs(G-044): set status=done`). The merge brought `approved` and the
verdict with the record. A merge that succeeds but a done write that fails
(the merged project does not validate, or Git refuses the commit) is
reported as exactly that: the merge stands, the record is not done, and the
command to finish is printed. Nothing undoes a merge.

`update` itself gains the one enforcement it could not have before:
where `grove.yaml` names a target, `--set status=done` is refused in a
checkout whose branch is not that target, naming both. Without a target
(nullsec today) the behaviour is unchanged. This closes G-074's first
finding with configuration the owner added in G-042.

Output is one line per fact, in order, on stdout as each completes:
`approval:` (the candidate, who approved it as the commit, the verdict),
`merge:` (fast-forward or merge commit and its hash, or the refusal),
`done:` (the commit), `cleanup:` (what was removed or kept). A refusal goes
to stderr with exit 1; the facts already printed stand. Approval and
integration are therefore always separate facts, as G-044 acceptance 4
asks.

### 4. Cleanup: opt-in, after done, by Git's own refusals

`integrate --cleanup` (and the board's second confirmation) removes the
branch's worktree and then the branch only after `done` is committed on the
target: `git worktree remove PATH` without `--force`, which refuses a
worktree holding changes or untracked files, then `git branch -d BRANCH`,
which refuses a branch the target does not contain. Grove adds one refusal
of its own: a worktree that holds the running process's working directory
is kept. A refusal is a `cleanup: kept …` fact with Git's reason and exit 1;
the integration stands. No branch or worktree is ever touched by `approve`,
`feedback`, a refused merge, or a failed done write. Records are never
moved (G-064).

Not taken: cleaning up by default (a kept worktree is cheap, a lost one is
not), and a standalone `cleanup` command (the justfile's `clean-merged` and
Git cover a later sweep).

### The CLI

```text
grove [--project DIR] approve ID VERDICT
grove [--project DIR] feedback ID TEXT
grove [--project DIR] integrate ID [--cleanup]
```

`approve` and `feedback` run in a checkout of the branch that holds the
record in `review` (the worktree an implementation used), because that is
where the record's review state lives. Both refuse a record that is not in
`review`, a checkout whose HEAD does not contain the candidate, a record
whose live file differs from HEAD (uncommitted changes are never committed
by these commands), and, for `approve`, a tip that differs from the
candidate in any file but the record (a new candidate). `approve` sets
`approved` to the candidate and appends the verdict; `feedback` sets
`status=active`, unsets `approved` when present, keeps the candidate (as
G-038 chose: the earlier review's `examined` still compares to it) and
appends the text. Both print what `update` prints, and `feedback` adds a
`Next:` line on stderr naming the branch, its checkout and `/grove-work ID`
there: the actionable interactive continuation of acceptance 3.

`integrate` runs in the target's checkout, as described above. It finds the
branch through the same inspection `versions` does, in one Git process.

### The board

The detail of a work record in `review` (G-096's review view) gains:

- **A Review block** in the header, facts only: the candidate; `approved` or
  `not approved`; whether the branch tip is the candidate (only the record
  changed since it) or names the commits after it; whether the candidate is
  on the target; and where the actions would run (the branch's checkout and
  the target's). Each linked review already shows what it examined against
  the candidate.
- **Content opens at Evidence** when the body has that heading, so the
  first screen is the handoff: outcome, changed behaviour, decisions,
  verification, findings.
- **A Changes section** in the sidebar, between Linked and Timeline, for any
  work with a candidate: the files changed from the merge base with the
  target to the candidate, with added and removed line counts, read on
  demand like the history (never during the board load, cancelled by any
  key) through two Git reads (`merge-base`, `diff --numstat`) plus one
  `diff --name-only` for the tip. Enter on a file shows its diff in the
  Content pane, every line escaped as record text is, with `+` and `-` rows
  coloured as an accent; Esc returns. That is the optional diff of
  acceptance 1.
- **Three actions**, each explicit and confirmed, on a review-status work
  record: `a` approve, which asks for the verdict on one input line (Enter
  runs, Esc cancels, empty is refused); `f` feedback, the same with the
  text; `i` integrate, which confirms `Merge BRANCH into TARGET in
  CHECKOUT and mark ID done? y/n`, then `Remove its worktree and branch?
  y/n`. The board runs the same functions the CLI does, in the checkout the
  Review block names, and shows the facts on a result screen (Esc back),
  then re-reads the board. While an action runs the banner says so and keys
  wait; an action's Git commands are never cancelled by a key. `a` and `f`
  need a clean checkout of the branch; `i` a checkout of the target; the
  board says which is missing rather than acting elsewhere.

Everything G-043 fixed stays: `safe` before any style, glamour only for
record bodies, focus and paging, on-demand history, one read at a time,
the explicit subcommands' contracts.

## Steps

1. **`approved`.** `internal/project`: the field, its two rules, `check`
   diagnostics; `internal/update`: the field tables, the drift map, and the
   done-off-target refusal (design 3); tests for each refusal.
2. **Approve and feedback.** `internal/update/review.go`: the body append
   (`Request.Append`), `Approve` and `Feedback` with their preconditions;
   `internal/cli`: the two commands and usage; tests in a Git fixture with a
   branch and a worktree: success, not in review, uncommitted record,
   commits after the candidate, checkout without the candidate, feedback
   after approval, the earlier review record untouched.
3. **Integrate.** `internal/integrate`: the checks of design 2, the merge,
   the done write of design 3, cleanup of design 4, and the fact lines;
   `internal/cli`: the command and `--cleanup`. Tests: fast-forward,
   merge commit after the target moved, conflict aborted with the target
   unchanged, the record edited on the target (a conflict), stale approval
   (a commit after the candidate), dirty target refused, wrong branch
   refused, no target refused, two branches refused, unapproved refused,
   cleanup done, cleanup kept (dirty worktree; the process's directory).
4. **Changes and diffs.** `internal/versions`: `Changes` and `Diff` reads
   with cancellation, tested on a fixture.
5. **Board.** Review block, Evidence start, Changes section and diff view,
   the prompt line, the three actions and the result screen, the checkout
   choice; tests with the fake backend (facts shown, prompt flow, refusals
   as notices, busy while acting, hostile diff text inert), the pty scenario
   extended through a review detail.
6. **Docs.** README (commands, the review view and its keys), the record
   model (`approved`, the lifecycle bullets, the "no approval field" note),
   `docs/work-execution.md` (step 8's integrator commands and the Judging
   section become the three commands), AGENTS.md's lifecycle sentence, the
   terms G-059 and G-060, and this plan's adjustments.
7. **Verification.** Per package while iterating; then `go vet ./...`,
   `gofmt -l .`, `go run ./cmd/grove check`, `go test -count=1 -timeout
   120s ./...`, `-race` for `internal/tui` and `internal/integrate` alone;
   a disposable clone exercised through every acceptance-4 path with the
   built binary; a real-terminal pass of the review detail, changes, a diff,
   and each action on the clone.
8. **Review and handoff.** An independent review of the final diff as a
   review record with `examined`; then G-044's Evidence and Next, and
   `status=review` with the candidate. The owner's judgment (acceptance 5)
   is the first real use: reviewing and integrating G-044 itself with
   `approve` and `integrate` from the board.

## Open choices

Put to the owner with this plan on 2026-09-22, each with the recommendation
above first, and answered the same day: the owner selected the recommended
option for all four.

1. Approval as an `approved` field bound to the candidate, with the verdict
   appended to the body (design 1), or verdict-at-integration only.
   **The field.**
2. A plain merge in the target's checkout, aborted on conflict (design 2),
   or fast-forward only. **Plain merge.**
3. `integrate` writes done on the target and `update` refuses done off a
   configured target (design 3), or `update` unchanged. **Both.**
4. Cleanup opt-in after done, by Git's own refusals (design 4), by default,
   or never from Grove. **Opt-in.**

## Adjustments

Made while implementing, none changing the four designs:

- The board's "pty scenario" is the connected-workflow test that drives the
  real model against real Git (`internal/cli/tui_test.go`), now built from
  `tui.Live()`, so the review detail, a diff, approval and integration run
  through the same backend the terminal uses.
- `feedback` prints its `Next:` continuation on stderr, and the board's
  result screen shows it as a fact line; there is no runner to hand it to.
- `update` gained `Branch(root)` for the done-off-target refusal, and
  `approve` runs the same tip check the guide describes so a checkout ahead
  of the candidate is refused with the new candidate named.
- `integrate` reads the process's working directory to keep a worktree it is
  running from; the board passes the same, so `i` from inside the branch's
  worktree keeps that worktree with the reason shown.

## Limits

Local integration only: no push, PR or deployment. Autonomous judging has no
policy here; every action is a person's key or command. Automatic relaunch
after feedback is G-046's. A branch without a checkout cannot be approved
from the board; `git worktree add` one. Diffs are per file from the merge
base; a whole-branch diff is `git diff` at a shell. Bodies other than the
appended paragraph are still edited by hand.
