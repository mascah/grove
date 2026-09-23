---
id: "G-044"
type: work
title: "Review candidates and integrate approved work locally"
status: active
created: "2026-09-21T00:54:16Z"
updated: "2026-09-23T00:25:14Z"
kind: feature
size: large
priority: 3
depends_on: ["G-038", "G-043"]
relates_to: ["G-035", "G-046", "G-064"]
formerly: "W-026"
---

## Outcome

Review an implementation in Grove, record feedback or approval against its
candidate, and integrate approved work locally with an explicit result.

## Scope and bounds

Lead with outcome, changed behavior, consequential decisions, verification,
open findings and follow-ups; progressively disclose artifacts, changed files
and optional diffs. Reuse G-038's review evidence rather than creating another
editable final report. Provide equivalent CLI operations for supported actions.

Bind approval to the candidate and relevant inputs/target state. Revalidate
before integration; conflicts, changed candidates and moved targets have clear
outcomes. Preserve partial work and dirty worktrees. Only clean up branches or
worktrees after integration is proven and retained files/evidence are safe.
Local integration only; no remote PR/push/deployment requirement.

Cleanup means safe branch/worktree cleanup, not moving completed records.
G-064's stable identity and placement apply to integration and feedback too.

## Acceptance

1. The owner can understand and judge a completed interactive candidate without
   the originating chat, inspect evidence, and optionally read diffs.
2. Approval, rejection/feedback and integration results remain attributable;
   a changed candidate cannot silently inherit approval.
3. Feedback requesting implementation returns work to Active and preserves the
   earlier review. Before a runner exists, it produces an actionable interactive
   continuation; automatic relaunch is added in G-046.
4. Exercise local merge success, conflict/refusal, changed target, stale approval,
   uncommitted changes and cleanup refusal. Integration failure does not mark
   Done or discard work. Report approval and integration as separate facts.
5. The owner accepts the review hierarchy and connected terminal workflow.

## Evidence

Branch `worktree-G-044` in `.claude/worktrees/G-044`, base `main` c4b3aea,
plan [G-098](G-098-g-044-review-integration-plan.md) (b17784d; the owner's
four choices recorded in 35992ac, which also set this record active at
revision `sha256:8577d1b1…`). Review [G-099](G-099-g-044-review-integration-review.md).
The candidate is the commit that records this evidence; `candidate` names
it. Commits: 11b5e37 `approved` and done-on-target; 2b899f1 approve and
feedback; 1bb3af8 integrate; cc6335b changes and diffs; d666dc3 and 660ac99
the board's review detail and a pseudo-terminal scenario; 0728074 docs;
2e0aca3, 956da77 and a0b2f2b the fixes from the exercise and the review.

Against the acceptance:

1. **Judging without the chat.** A record in review opens on the board with
   a Review block (candidate, approved or not, whether only the record
   changed since it, whether the target holds it, where `a`/`f`/`i` run),
   its content at `## Evidence`, the linked reviews with their `examined`
   against the candidate, and the changed files against the target with
   counts; Enter on a file shows its diff, escaped. `context G-044` and
   `show` remain the noninteractive route.
2. **Attributable, bound to the candidate.** `approve` sets the work field
   `approved`, valid only when equal to `candidate` and in review or done,
   and appends `Verdict on candidate X, DATE: …`; a changed candidate fails
   `check`, `approve` refuses a tip that changed other files, `integrate`
   re-checks the tip and merges the inspected commit. Every write is a
   commit of the record alone with a generated message.
3. **Feedback.** `feedback` appends `Feedback on candidate X, DATE: …`,
   sets `active`, unsets `approved`, keeps `candidate` so G-099-style
   reviews still compare, and prints `Next: … continue there with
   /grove-work ID`; the board shows the same as a fact.
4. **Integration paths.** Exercised on a disposable clone with the built
   binary and in `internal/integrate` tests: fast-forward; merge commit
   after the target moved; conflict aborted with the target unchanged and
   the record left in review (the refusal names the conflicting file);
   stale approval (a later commit changes code) refused; uncommitted record
   refused by `approve` and `feedback`; dirty target refused; unapproved,
   wrong branch, no target and two branches refused; `--cleanup` keeping a
   worktree with untracked, ignored or the process's own directory, and a
   branch Git will not delete, with the integration standing; a failed done
   write reported after the merge stands. Approval, merge, done and cleanup
   are printed as separate fact lines.
5. **The owner's judgment** is Next: reviewing and integrating this record
   with the board and the commands is the first real use.

Verification at a0b2f2b, whose code the candidate shares: `go vet ./...`
and `gofmt -l .` clean; `go run ./cmd/grove check` OK, 95 records;
`go test -count=1 -timeout 120s ./...` passes (`TestTerminal` runs nine
pseudo-terminal scenarios, the new one three more times alone); `-race`
alone on `internal/integrate` and `internal/tui` passes; `-short` alone:
versions 4.2 s, cli 2.2 s, tui 0.4 s, integrate 1.5 s, update 1.7 s (the
suite's parallel run shows versions at 9 s and cli at 6 s, as before this
work). The clone exercise (twenty-two commands) ran before and after the
review's fixes with the same outcomes.

Decisions: a plain merge of the inspected commit, aborted on conflict;
`integrate` writes done and `update` refuses done off a configured target;
cleanup opt-in through Git's refusals plus ignored files; the y/n prompts
lead with the question so a long path is what truncation drops; a
pseudo-terminal expectation is one draw's text, since the renderer scrolls
and redraws lines from their first changed cell.

Limits: local only, no push or PR; a branch without a checkout cannot be
judged from the board; Ctrl-C during an action waits and reports only the
interruption; `update --set approved=` by hand skips the verdict and tip
check; diffs are per file from the merge base.

## Next

The owner judges this candidate, which is the first use of the review view.
`main`'s build lacks these commands, so build this branch once:

```sh
cd /Users/mascah/GitHub/mascah/grove/.claude/worktrees/G-044
go build -o /tmp/grove-G-044 ./cmd/grove
```

Then either open the board from a clean checkout of `main` and press Enter
on G-044, `a` for the verdict and `i` to merge (answer the cleanup prompt
`y` to remove this worktree and branch):

```sh
cd /Users/mascah/GitHub/mascah/grove && /tmp/grove-G-044
```

or run the commands, `approve` in this worktree and `integrate` in `main`'s
checkout, whose tracked files must be clean (it had an uncommitted
`internal/tui/testdata/terminal.py` edit on 2026-09-23):

```sh
/tmp/grove-G-044 --project /Users/mascah/GitHub/mascah/grove/.claude/worktrees/G-044 approve G-044 "VERDICT"
cd /Users/mascah/GitHub/mascah/grove && /tmp/grove-G-044 integrate G-044 --cleanup
```

Feedback instead: `/tmp/grove-G-044 --project …/G-044 feedback G-044 "TEXT"`
or `f` on the board, then `/grove-work G-044` in this worktree.
